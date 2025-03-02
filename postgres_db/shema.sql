CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(20) UNIQUE,
    address TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID REFERENCES customers(id) NOT NULL,
    account_number VARCHAR(20) UNIQUE NOT NULL,
    account_type VARCHAR(50) NOT NULL, -- Ex - "checking", "savings", "credit"
    balance DECIMAL(19, 4) NOT NULL DEFAULT 0.0000,
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- Ex- "active", "closed", "frozen"
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 0
);

-- Indexes for performance
CREATE INDEX idx_accounts_customer_id ON accounts (customer_id);
CREATE INDEX idx_accounts_account_number ON accounts (account_number);
CREATE INDEX idx_customers_email ON customers (email);
CREATE INDEX idx_customers_phone_number ON customers (phone_number);


-- Function to automatically update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to call the function before updating accounts
CREATE TRIGGER update_accounts_updated_at
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Trigger to call the function before updating customers
CREATE TRIGGER update_customers_updated_at
BEFORE UPDATE ON customers
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

