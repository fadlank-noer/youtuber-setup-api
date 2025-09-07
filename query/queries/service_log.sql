-- name: CreateServiceLog :exec
INSERT INTO service_logs (
  service_name, metadata
) VALUES (
  $1, $2
);