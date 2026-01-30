package models

type EmailAttachment struct {
	PdfBase64 string `json:"pdf_base64"`
	UrlFileUp string `json:"urlFileUp"`
}

type PayloadVerificacion struct {
	FileUp    string `json:"fileUp"`
	UrlFileUp string `json:"urlFileUp"`
}
