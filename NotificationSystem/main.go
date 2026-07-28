package main

import "fmt"

func main(){
	var notifiers = []Notifier{
		&EmailNotifier{toAddress: "giridhar.gabbur@gmail.com"},
		&SMSNotifier{phoneNumber: "7899455667"},
		&SlackNotifier{channel: "teams"},
	}

	NotifyAll(notifiers, "Hiya how you doin")
	NotifyFirst(notifiers, "Heya")
	retryemail := RetryNotifier{
		notifier: &EmailNotifier{toAddress: "giridhar.gabbur@gmail.com"},
		maxRetries: 3,
	}
	retryemail.send("Hey man wassup")
}

type Notifier interface {
	Send(message string) error
	Name() string
}

type EmailNotifier struct {
	toAddress string
}

type SMSNotifier struct {
	phoneNumber string
}

type SlackNotifier struct {
	channel string
}

type RetryNotifier struct{
	notifier Notifier
	maxRetries int
}

func (e *EmailNotifier) Send(message string) error {
	if message == "" {
		return fmt.Errorf("Message cannot be empty")
	}
	fmt.Printf("Sending email to %s: %s\n", e.toAddress,message)
	return nil
}

func (e *EmailNotifier) Name() string{
	return "This is an Email notifer"
}

func (s *SMSNotifier) Send(message string) error {
	if message == "" {
		return fmt.Errorf("Message cannot be empty")
	}
	fmt.Printf("Sending SMS to %s: %s\n", s.phoneNumber,message)
	return nil
}

func (s *SMSNotifier) Name() string{
	return "This is an SMS notifer"
}

func (k *SlackNotifier) Send(message string) error {
	if message == "" {
		return fmt.Errorf("Message cannot be empty")
	}
	fmt.Printf("Sending slack message to %s: %s\n", k.channel,message)
	return nil
}

func (k *SlackNotifier) Name() string{
	return "This is an Slack notifer"
}

func NotifyAll (notifiers []Notifier, message string){
	for _,n := range notifiers {
		fmt.Println(n.Name())
		err := n.Send(message)
		if err != nil{
			fmt.Println("Messsage failed to send")
		}
	}
}

func NotifyFirst (notifiers []Notifier, message string) error{
	for _,n := range notifiers {
		fmt.Println(n.Name())
		err := n.Send(message)
		if err != nil {
			fmt.Printf("%s Failed, trying next", n.Name())
			continue
		}
		fmt.Printf("%s succeeded!\n", n.Name())
		return nil
	}
	return fmt.Errorf("All notifiers failed")
}

func (r *RetryNotifier) send(message string) error{
	for attempt := 0; attempt <= r.maxRetries; attempt++{
		err := r.notifier.Send(message)
		if (err == nil){
			fmt.Printf("Succeeded on attempt: %d\n",attempt)
			return nil
		}
		fmt.Printf("Attempt %d/%d failed: %v, retrying",attempt, r.maxRetries, err)
	}
	return fmt.Errorf("%s failed after %d attempts", r.notifier.Name(), r.maxRetries)
}