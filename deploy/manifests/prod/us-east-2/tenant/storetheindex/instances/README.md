# Instances

List of individually configurable instances:

| Instance | sth bit-size | IOPS per GiB  | Value Codec  | Whitelist           | `STHBurstRate` | `STHSyncInterval` | Running                                                                                                                                       |
|----------|--------------|---------------|--------------|---------------------|----------------|-------------------|-----------------------------------------------------------------------------------------------------------------------------------------------|
| `romi`   | 30           | 5             | `json`       | all                 | `4194304`      | `1s`              | [778339d270108841997806c86203ddd3a7341fcb](https://github.com/filecoin-project/storetheindex/commit/778339d270108841997806c86203ddd3a7341fcb) |
| `tara`   | 30           | 5             | `json`       | all                 | `83886080`     | `6s`              | [778339d270108841997806c86203ddd3a7341fcb](https://github.com/filecoin-project/storetheindex/commit/778339d270108841997806c86203ddd3a7341fcb) |
| `xabi`   | 30           | 5             | `binary`     | all                 | `4194304`      | `1s`              | [945940507682064093e846ecc8578a58a5f16535](https://github.com/filecoin-project/storetheindex/commit/945940507682064093e846ecc8578a58a5f16535) |
| `vega`   | 30           | 5             | `binary`     | `nft.storage` only  | `4194304`      | `3s`              | [778339d270108841997806c86203ddd3a7341fcb](https://github.com/filecoin-project/storetheindex/commit/778339d270108841997806c86203ddd3a7341fcb) |
| `oden`   | N/A Pebble   | 5             | N/A Pebble   | `nft.storage` only  | N/A Pebble     | N/A Pebble        | [498f406a73629637654c936c5d34ab768f2e417d](https://github.com/filecoin-project/storetheindex/commit/498f406a73629637654c936c5d34ab768f2e417d) |
| `dido`   | N/A Pebble   | 5             | N/A Pebble   | all                 | N/A Pebble     | N/A Pebble        | [1705a46c927043fb3e2c2d34de4855bf39c906b3](https://github.com/filecoin-project/storetheindex/commit/1705a46c927043fb3e2c2d34de4855bf39c906b3) |
