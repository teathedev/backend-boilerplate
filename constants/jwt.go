package constants

// JWTMasterKeyEnv is the environment variable for the AES key used to encrypt/decrypt
// JWT private keys in storage. Must be base64-encoded 32 bytes for AES-256-GCM.
const JWTMasterKeyEnv = "JWT_MASTER_KEY"
