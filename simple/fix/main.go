package main

import (
	"context"
	"fmt"

	clientset "github.com/tilt-dev/tilt-api-client-go/pkg/clientset/versioned"
	"github.com/tilt-dev/tilt-api-client-go/pkg/config"
	"github.com/tilt-dev/tilt/pkg/apis/core/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	tiltAPIConfig, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	cli := clientset.NewForConfigOrDie(tiltAPIConfig)
	w, err := cli.TiltV1alpha1().Sessions().Watch(context.Background(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// // everything on every change
	// for e := range w.ResultChan() {
	// 	fmt.Println(e.Object)
	// }

	// // diff
	// for e := range w.ResultChan() {
	// 	// for _, t := range e.Object.(*v1alpha1.Session).Status.Targets {
	// 	//   if t.State.Terminated != nil && t.State.Terminated.Error != "" {
	// 	fmt.Println(e.Object) // t.Name
	// 	//   }
	// 	// }
	// }

	// // only errors
	// for e := range w.ResultChan() {
	// 	for _, t := range e.Object.(*v1alpha1.Session).Status.Targets {
	// 		if t.State.Terminated != nil && t.State.Terminated.Error != "" {
	// 			// fmt.Println(t.Name)
	// 			revolutionizeoutsidethebox(t.Name)
	// 		}
	// 	}
	// }

	// // diff
	// //	prevFails := make(map[string]bool)
	// for e := range w.ResultChan() {
	// 	// fails := make(map[string]bool)
	// 	for _, t := range e.Object.(*v1alpha1.Session).Status.Targets {
	// 		if t.State.Terminated != nil && t.State.Terminated.Error != "" {
	// 			// fails[t.Name] = true
	// 			// if !prevFails[t.Name] {
	// 			fmt.Println(t.Name)
	// 			// revolutionizeoutsidethebox(t.Name)
	// 		}
	// 	}
	// }
	// // prevFails = fails
	// // }

	// diff
	prevFails := make(map[string]bool)
	for e := range w.ResultChan() {
		fails := make(map[string]bool)
		for _, t := range e.Object.(*v1alpha1.Session).Status.Targets {
			if t.State.Terminated != nil && t.State.Terminated.Error != "" {
				fails[t.Name] = true
				if !prevFails[t.Name] {
					fmt.Println(t.Name)
					revolutionizeoutsidethebox(t.Name)
				}
			}
		}
		prevFails = fails
	}
}
