package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DamnWidget/goqueue"
)

// Queue for incoming alerts
var alert_queue *goqueue.Queue

// This struct will be pushed onto the queue
type alert_entry struct {
	time_stamp int64  // Timestamp of alert
	alert_msg  string // Alert message
}

// Listening interface
const iface_ip string = "0.0.0.0"

// Service port
const service_port string = "8080"

// Store an alert in the queue
func set_alert(writer http.ResponseWriter, request *http.Request) {
	// This gets a map of keys and values
	alert_all_values := request.URL.Query()

	if request.Method == http.MethodPost {
		request.ParseForm()
	}
	alert_post_values := request.Form

	if alert_all_values == nil {
		alert_queue.Push(&alert_entry{time.Now().Unix(), "GET EMPTY"})
	} else {
		alert_queue.Push(&alert_entry{time.Now().Unix(), fmt.Sprint(alert_all_values)})
	}

	if alert_post_values == nil {
		alert_queue.Push(&alert_entry{time.Now().Unix(), "POST EMPTY"})
	} else {
		alert_queue.Push(&alert_entry{time.Now().Unix(), fmt.Sprint(alert_post_values)})
	}

	// // Push a string version of the map to the queue with a timestamp
	// if alert_all_values != nil {
	// alert_queue.Push(&alert_entry{time.Now().Unix(), fmt.Sprint(alert_all_values)})
	// alert_queue.Push(&alert_entry{time.Now().Unix(), fmt.Sprint(alert_post_values)})
	// }
}

// Send the queued alerts to the client
func get_alert(writer http.ResponseWriter, request *http.Request) {
	if alert_queue.Len() > 0 {
		for {
			entry := alert_queue.Pop()
			if entry == nil {
				break
			}

			// See if the entry can be cast to the proper type
			alert_entrty, valid := entry.(*alert_entry)
			if !valid {
				panic("invalid type")
			}

			// Write response to the output stream
			fmt.Fprintln(writer, strconv.FormatInt(alert_entrty.time_stamp, 10)+": "+alert_entrty.alert_msg)
		}
	}
}

// Entry point
func main() {
	// Create new queue, this is thread safe
	alert_queue = goqueue.New()

	// set up the route handlers
	http.HandleFunc("/api/alert/set", set_alert)
	http.HandleFunc("/api/alert/get", get_alert)

	// Listen on the specified interface and port
	// TODO: Make these env variables
	log.Fatal(http.ListenAndServe(iface_ip+":"+service_port, nil))
}
