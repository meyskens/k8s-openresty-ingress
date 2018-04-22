package connector

import (
	"errors"
	"fmt"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type watchable interface {
	Watch(meta_v1.ListOptions) (watch.Interface, error)
}

func (c *Client) watcher(changeChan chan bool, object watchable) error {
	w, err := object.Watch(meta_v1.ListOptions{})
	if err != nil {
		return err
	}

	for {
		event := <-w.ResultChan()
		if event.Type == watch.Error || event.Object == nil {
			err = errors.New(fmt.Sprintln(event.Object))
			break
		}
		changeChan <- true
	}
	w.Stop()

	return err
}
