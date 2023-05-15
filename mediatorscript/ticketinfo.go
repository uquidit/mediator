package mediatorscript

import (
	"encoding/xml"
)

type TicketInfo struct {
	XMLName          xml.Name     `xml:"ticket_info"`
	ID               int          `xml:"id"`
	Subject          string       `xml:"subject"`
	PriorityID       int          `xml:"priority>id"`
	PriorityName     string       `xml:"priority>name"`
	CreateDate       int          `xml:"createDate"`
	UpdateDate       int          `xml:"updateDate"`
	Requester        TicketUser   `xml:"requester"`
	CurrentStage     *TicketStage `xml:"current_stage"`
	CompletionData   *TicketStage `xml:"completion_data>stage"`
	OpenRequestStage TicketStage  `xml:"open_request_stage"`
	Comment          string       `xml:"comment"`
}

type TicketStage struct {
	ID          int        `xml:"id"`
	Name        string     `xml:"name"`
	TaskID      int        `xml:"ticket_task>id"`
	TaskName    string     `xml:"ticket_task>name"`
	TaskHandler TicketUser `xml:"ticket_task>handler"`
}

type TicketUser struct {
	ID          int    `xml:"id"`
	Login       string `xml:"login"`
	DisplayName string `xml:"display_name"`
}
