output "image_id" {
  description = "The ID of the Docker image"
  value       = docker_registry_image.this.id
}