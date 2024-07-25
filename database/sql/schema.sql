CREATE TABLE Member_Roles (
    id SERIAL PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE Transaction_Types (
    id SERIAL PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE Event_Types (
    id SERIAL PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE Users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc')
);

CREATE TABLE Accounts (
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    account_name TEXT NOT NULL,
    account_type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT fk_account_user FOREIGN KEY (user_id) REFERENCES Users(id)
);

CREATE TABLE Account_Events (
    id SERIAL PRIMARY KEY,
    account_id INT NOT NULL,
    event_type_id INT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT fk_account_event_account FOREIGN KEY (account_id) REFERENCES Accounts(id),
    CONSTRAINT fk_account_event_event_type FOREIGN KEY (event_type_id) REFERENCES Event_Types(id)
);

CREATE TABLE Members (
    id SERIAL PRIMARY KEY,
    account_id int NOT NULL,
    user_id int NOT NULL,
    member_role_id int NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    CONSTRAINT fk_member_account FOREIGN KEY (account_id) REFERENCES Accounts(id),
    CONSTRAINT fk_member_user FOREIGN KEY (user_id) REFERENCES Users(id),
    CONSTRAINT fk_member_member_role FOREIGN KEY (member_role_id) REFERENCES Member_Roles(id)
);

CREATE TABLE Transactions (
    id SERIAL PRIMARY KEY,
    account_id int NOT NULL,
    user_id int NOT NULL,
    transaction_date DATE NOT NULL,
    transaction_type_id int NOT NULL,
    amount FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'utc'),
    description TEXT NOT NULL,
    CONSTRAINT fk_transaction_account FOREIGN KEY (account_id) REFERENCES Accounts(id),
    CONSTRAINT fk_transaction_user FOREIGN KEY (user_id) REFERENCES Users(id),
    CONSTRAINT fk_transaction_transaction_type FOREIGN KEY (transaction_type_id) REFERENCES Transaction_Types(id)
);