package astisparkpost

import (
	"fmt"
	"net/http"
)

type Transmission struct {
	CampaignID       string               `json:"campaign_id,omitempty"`
	Content          *InlineContent       `json:"content,omitempty"`
	Description      string               `json:"description,omitempty"`
	Metadata         *Metadata            `json:"metadata,omitempty"`
	Options          *TransmissionOptions `json:"options,omitempty"`
	Recipients       []Recipient          `json:"recipients,omitempty"`
	ReturnPath       string               `json:"return_path,omitempty"`
	SubstitutionData *SubstitutionData    `json:"substitution_data,omitempty"`
}

type TransmissionOptions struct {
	ClickTracking *bool `json:"click_tracking,omitempty"`
	Sandbox       *bool `json:"sandbox,omitempty"`
}

type Recipient struct {
	Address          *Address          `json:"address,omitempty"`
	Metadata         *Metadata         `json:"metadata,omitempty"`
	ReturnPath       string            `json:"return_path,omitempty"`
	SubstitutionData *SubstitutionData `json:"substitution_data,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
}

type Address struct {
	Email    string `json:"email,omitempty"`
	HeaderTo string `json:"header_to,omitempty"`
	Name     string `json:"name,omitempty"`
}

type InlineContent struct {
	Attachments  []Attachment  `json:"attachments,omitempty"`
	From         *Address      `json:"from,omitempty"`
	Headers      *Headers      `json:"headers,omitempty"`
	HTML         string        `json:"html,omitempty"`
	InlineImages []InlineImage `json:"inline_images,omitempty"`
	ReplyTo      string        `json:"reply_to,omitempty"`
	Subject      string        `json:"subject,omitempty"`
	TemplateID   string        `json:"template_id,omitempty"`
	Text         string        `json:"text,omitempty"`
}

type Headers map[string]string

type Metadata map[string]string

type SubstitutionData map[string]interface{}

type Attachment struct {
	Data string `json:"data,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type InlineImage struct {
	Data string `json:"data,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type TransmissionResults struct {
	ID                      string `json:"id"`
	TotalAcceptedRecipients int    `json:"total_accepted_recipients"`
	TotalRejectedRecipients int    `json:"total_rejected_recipients"`
}

// CreateLink creates a link
func (c *Client) CreateTransmission(t Transmission) (r TransmissionResults, err error) {
	// Send
	if err = c.send(http.MethodPost, "/v1/transmissions", t, &r); err != nil {
		err = fmt.Errorf("astisparkpost: sending failed: %w", err)
		return
	}
	return
}
