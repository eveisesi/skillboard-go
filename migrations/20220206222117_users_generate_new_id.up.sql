UPDATE
    `users`
SET
    id = SHA1(CONCAT(owner_hash, ":", character_id));