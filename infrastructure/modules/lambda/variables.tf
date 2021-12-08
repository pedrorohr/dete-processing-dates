variable "name" {
    type        = string 
    description = "A name to define the lambda"
}

variable "handler" {
    type        = string 
    description = "Entrypoint of the lambda function"
}

variable "source_file" {
    type        = string 
    description = "Path to the lambda source file"
}

variable "extra_policies" {
    type        = map
    description = "Extra policies to be attached to the lambda role"
    default     = {}
}

variable "env" {
    type        = map
    description = "Environment variables accessible from the lambda function during execution"
    default     = {}
}
