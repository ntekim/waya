CREATE TABLE payouts (
    id TEXT PRIMARY KEY,
    batch_id TEXT,
    reference_id TEXT NOT NULL, -- Client's unique ID for this payout
    
    recipient_name TEXT NOT NULL,
    recipient_phone TEXT NOT NULL,
    recipient_email TEXT, -- Can be NULL
    recipient_tag TEXT,   -- Can be NULL
    country_code TEXT NOT NULL,

    -- Bank Details (for bank transfers)
    bank_code TEXT,       -- 033 (UBA) -- Can be NULL (if mobile money)
    bank_name TEXT,
    account_number TEXT,  -- 1234567890
    
    -- Money
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    
    status TEXT NOT NULL,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE batches (
    id TEXT PRIMARY KEY,
    total_amount BIGINT NOT NULL,
    total_count INTEGER NOT NULL,
    status TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);