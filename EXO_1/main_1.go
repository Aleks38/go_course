package main

import "fmt"

type Notifier interface {
    Send(message string) error
}

type EmailNotifier struct {
    Recipient string
    Sender    string
}

func (e EmailNotifier) Send(message string) error {
    fmt.Printf("[EMAIL] De %s à %s : %s\n", e.Sender, e.Recipient, message)
    return nil
}

type SMSNotifier struct {
    PhoneNumber string
}

func (s SMSNotifier) Send(message string) error {
    fmt.Printf("[SMS] Envoi à %s : %s\n", s.PhoneNumber, message)
    return nil
}

type ConsoleNotifier struct{}

func (c ConsoleNotifier) Send(message string) error {
    fmt.Printf("[CONSOLE] Message : %s\n", message)
    return nil
}

func main() {
    notifiers := []Notifier{
        EmailNotifier{Recipient: "alice@example.com", Sender: "bob@example.com"},
        SMSNotifier{PhoneNumber: "+33123456789"},
        ConsoleNotifier{},
    }

    for _, n := range notifiers {
        if err := n.Send("Bonjour depuis Go !"); err != nil {
            fmt.Println("Erreur lors de l'envoi :", err)
        }
    }
}

