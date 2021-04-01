package config

// nolint:lll
const defaultConfig = `
server:
  address: :65432
  read-timeout: 20s
  write-timeout: 20s
  graceful-timeout: 5s
jwt:
  expiration: 2h
  secret: 'SECRET_TOKEN'
`
