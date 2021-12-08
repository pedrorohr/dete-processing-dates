output "arn" {
  description = "ARN of the lambda"
  value       = aws_lambda_function.lambda.arn
}

output "function_name" {
  description = "Function name of the lambda"
  value       = aws_lambda_function.lambda.function_name 
}