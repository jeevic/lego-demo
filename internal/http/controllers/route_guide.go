package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jeevic/lego/components/grpc/grpcclient"

	routeGuide "github.com/jeevic/lego-demo/pb/route_guide"
)

func GetFeature(c *gin.Context) {
	p := &routeGuide.Point{Latitude: 409146138, Longitude: -746188906}

	str := fmt.Sprintf("Getting feature for point (%d, %d) \n", p.Latitude, p.Longitude)

	addr := "127.0.0.1:8501"
	conn, _ := grpcclient.NewClient(addr, grpcclient.NewOptions())
	defer conn.Close()
	cli := routeGuide.NewRouteGuideClient(conn.GetConn())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	feature, err := cli.GetFeature(ctx, p)
	if err != nil {
		str = str + fmt.Sprintf("%v.GetFeatures(_) = _, %v: ", cli, err)
	}
	str = str + fmt.Sprintf("feature:%#v", feature)

	c.String(http.StatusOK, str)
}

func GetFeatures(c *gin.Context) {
	// Looking for features between 40, -75 and 42, -73.
	rect := &routeGuide.Rectangle{
		Lo: &routeGuide.Point{Latitude: 400000000, Longitude: -750000000},
		Hi: &routeGuide.Point{Latitude: 420000000, Longitude: -730000000},
	}

	addr := "127.0.0.1:8501"
	conn, _ := grpcclient.NewClient(addr, grpcclient.NewOptions())
	defer conn.Close()
	cli := routeGuide.NewRouteGuideClient(conn.GetConn())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	str := ""
	stream, err := cli.ListFeatures(ctx, rect)
	if err != nil {
		str = fmt.Sprintf("%v.ListFeatures(_) = _, %v \n", cli, err)
	}

	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			str = str + "this is io EOF \n"
			break
		}
		if err != nil {
			str = str + fmt.Sprintf("%v.ListFeatures(_) = _, %v", cli, err)
		}
		str = str + fmt.Sprintf("Feature: name: %q, point:(%v, %v) \n", feature.GetName(),
			feature.GetLocation().GetLatitude(), feature.GetLocation().GetLongitude())
	}

	stream, err = cli.ListFeatures(ctx, rect)
	if err != nil {
		str = fmt.Sprintf("%v.ListFeatures(_) = _, %v \n", cli, err)
	}

	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			str = str + "this is io EOF \n"
			break
		}
		if err != nil {
			str = str + fmt.Sprintf("%v.ListFeatures(_) = _, %v", cli, err)
		}
		str = str + fmt.Sprintf("Feature: name: %q, point:(%v, %v) \n", feature.GetName(),
			feature.GetLocation().GetLatitude(), feature.GetLocation().GetLongitude())
	}

	c.String(http.StatusOK, str)
}

func RunRecordRoute(c *gin.Context) {
	addr := "127.0.0.1:8501"
	conn, _ := grpcclient.NewClient(addr, grpcclient.NewOptions())
	defer conn.Close()
	cli := routeGuide.NewRouteGuideClient(conn.GetConn())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a random number of random points
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
	var points []*routeGuide.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}

	str := fmt.Sprintf("Traversing %d points.", len(points))

	stream, err := cli.RecordRoute(ctx)
	if err != nil {
		str = fmt.Sprintf("%v.RecordRoute(_) = _, %v", cli, err)
		c.String(http.StatusOK, str)
		return
	}
	header, err := stream.Header()
	if err != nil {
		str = fmt.Sprintf("metadata:%v", header)
	}

	for _, point := range points {
		if err := stream.Send(point); err != nil {
			str = str + fmt.Sprintf("%v.Send(%v) = %v", stream, point, err)
			c.String(http.StatusOK, str)
			return
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		str = str + fmt.Sprintf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
		c.String(http.StatusOK, str)
		return
	}
	str = str + fmt.Sprintf("Route summary: %v", reply)
	c.String(http.StatusOK, str)
	return
}

func RunRouteChat(c *gin.Context) {
	addr := "127.0.0.1:8501"
	conn, _ := grpcclient.NewClient(addr, grpcclient.NewOptions())
	defer conn.Close()
	cli := routeGuide.NewRouteGuideClient(conn.GetConn())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	notes := []*routeGuide.RouteNote{
		{Location: &routeGuide.Point{Latitude: 0, Longitude: 1}, Message: "First message"},
		{Location: &routeGuide.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &routeGuide.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &routeGuide.Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
		{Location: &routeGuide.Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
		{Location: &routeGuide.Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
	}

	str := ""

	stream, err := cli.RouteChat(ctx)
	if err != nil {
		str = fmt.Sprintf("%v.RouteChat(_) = _, %v", cli, err)
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				str = str + "io EOF"
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				str = str + fmt.Sprintf("Failed to receive a note : %v", err)
				break
			}
			str = str + fmt.Sprintf("Got message %s at point(%d, %d) \n", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()

	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}

	_ = stream.CloseSend()
	<-waitc
	c.String(http.StatusOK, str)
}

func randomPoint(r *rand.Rand) *routeGuide.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &routeGuide.Point{Latitude: lat, Longitude: long}
}
