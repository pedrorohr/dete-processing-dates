variable "dete_processing_dates_url" {
  type        = string 
  description = "URL where DETE processing dates can be found"
}

variable "dete_bot_api_token" {
  type        = string
  description = "Telegram Bot API Token to be used to send DETE processing dates updates"
  sensitive   = true
}

variable "dete_chat_id" {
  type        = string 
  description = "Telegram Chat ID to send DETE processing dates updates"
  sensitive   = true
}
