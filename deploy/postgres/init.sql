-- Postgres initialization for Accounts, OperationTypes, and Transactions
CREATE TABLE IF NOT EXISTS Accounts (
    Id SERIAL PRIMARY KEY,
    DocumentNumber VARCHAR NOT NULL UNIQUE,
    CreatedAt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS OperationTypes (
    Id SERIAL PRIMARY KEY,
    Description VARCHAR NOT NULL
);

CREATE TABLE IF NOT EXISTS Transactions (
    Id SERIAL PRIMARY KEY,
    AccountId INT NOT NULL REFERENCES Accounts(Id),
    OperationTypeId INT NOT NULL REFERENCES OperationTypes(Id),
    Amount NUMERIC NOT NULL,
    EventDate TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Insert initial data into OperationTypes
INSERT INTO OperationTypes (Description) VALUES
('Normal Purchase'),
('Purchase with installments'),
('Withdrawal'),
('Credit Voucher');
