package config

// nolint:lll
const defaultConfig = `
server:
  address: :65432
  read-timeout: 2m
  write-timeout: 2m
  graceful-timeout: 5s
jwt:
  expiration: 2h
  secret: 'SECRET_TOKEN'
`
