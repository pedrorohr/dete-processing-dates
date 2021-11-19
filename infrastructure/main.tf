terraform {
  required_version = "~> 1.0.0"
}

data "archive_file" "scraper_zip" {
  type        = "zip"
  source_file = "../bin/scraper"
  output_path = "bin/scraper.zip"
}

