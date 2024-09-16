CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE organization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    type organization_type,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    user_id UUID REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TABLE organization_employee (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    username VARCHAR(50) REFERENCES employee(username) ON DELETE CASCADE
);


CREATE TYPE tender_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);

CREATE TYPE tender_status AS ENUM (
    'Created',
    'Published',
    'Closed'
);


CREATE TABLE tender (
    tender_id UUID NOT NULL, 
    tender_name VARCHAR(100) NOT NULL, 
    tender_description VARCHAR(100),
    service_type tender_type, 
    status tender_status, 
    organization_id UUID NOT NULL REFERENCES organization(id),
    tender_version INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tender_id, tender_version)
);


CREATE TYPE author_bid AS ENUM (
    'Organization',
    'User'
);

CREATE TYPE bid_status AS ENUM(
    'Created',
    'Published',
    'Canceled',
    'Rejected',
    'Approved'
);

CREATE TABLE bid (
    bid_id UUID NOT NULL,  
    bid_name VARCHAR(100) NOT NULL, 
    bid_description VARCHAR(500),
    bid_status bid_status, 
    tender_id UUID NOT NULL,
    tender_version INT NOT NULL DEFAULT 0, 
    approved_count INT DEFAULT 0, 
    bid_author_type author_bid,
    bid_author_id UUID NOT NULL REFERENCES employee(id),
    bid_version INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (bid_id, bid_version),  
    FOREIGN KEY (tender_id, tender_version) REFERENCES tender(tender_id, tender_version) ON DELETE CASCADE
);




CREATE TABLE bid_feedback (
    id UUID PRIMARY KEY NOT NULL,
    bid_id UUID NOT NULL,
    bid_version INT NOT NULL,
    author UUID REFERENCES employee(id),
    bid_feedback VARCHAR(1000),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (bid_id, bid_version) REFERENCES bid(bid_id, bid_version) ON DELETE CASCADE
);



CREATE TABLE bid_decision(
    id UUID PRIMARY KEY NOT NULL, 
    bid_id UUID NOT NULL,
    bid_version INT NOT NULL,
    author UUID REFERENCES employee(id),
    bid_feedback VARCHAR(1000),
    FOREIGN KEY (bid_id, bid_version) REFERENCES bid(bid_id, bid_version) ON DELETE CASCADE
);

-- Вставляем сотрудников
INSERT INTO employee (username, first_name, last_name) VALUES
('john_doe', 'John', 'Doe'),
('jane_smith', 'Jane', 'Smith'),
('alice_johnson', 'Alice', 'Johnson'),
('bob_brown', 'Bob', 'Brown'),
('darya_emelyanenko', 'Darya', 'Emelyanenko');

-- Вставляем организации
INSERT INTO organization (name, description, type) VALUES
('Tech Innovators LLC', 'A leading tech company specializing in innovative solutions.', 'LLC'),
('Green Energy IE', 'An organization focused on renewable energy solutions.', 'IE'),
('Global Manufacturing JSC', 'A multinational corporation in the manufacturing sector.', 'JSC');


-- Вставляем ответственных за организацию
INSERT INTO organization_responsible (organization_id, user_id) VALUES
((SELECT id FROM organization WHERE name = 'Tech Innovators LLC'), (SELECT id FROM employee WHERE username = 'john_doe')),
((SELECT id FROM organization WHERE name = 'Green Energy IE'), (SELECT id FROM employee WHERE username = 'jane_smith')),
((SELECT id FROM organization WHERE name = 'Global Manufacturing JSC'), (SELECT id FROM employee WHERE username = 'alice_johnson'));


INSERT INTO organization_responsible (organization_id, user_id)
VALUES (
    (SELECT id FROM organization WHERE name = 'Global Manufacturing JSC'),
    (SELECT id FROM employee WHERE username = 'bob_brown')
);

-- Вставляем сотрудников в организацию
INSERT INTO organization_employee (organization_id, username) VALUES
((SELECT id FROM organization WHERE name = 'Tech Innovators LLC'), 'john_doe'),
((SELECT id FROM organization WHERE name = 'Green Energy IE'), 'jane_smith'),
((SELECT id FROM organization WHERE name = 'Global Manufacturing JSC'), 'alice_johnson'),
((SELECT id FROM organization WHERE name = 'Global Manufacturing JSC'), 'bob_brown'),
((SELECT id FROM organization WHERE name = 'Global Manufacturing JSC'), 'darya_emelyanenko');


INSERT INTO employee (username, first_name, last_name) VALUES
('test_user', 'Test', 'User');
INSERT INTO organization (name, description, type) VALUES
('Test Organization', 'A sample organization for testing.', 'LLC');
INSERT INTO organization_responsible (organization_id, user_id) VALUES
((SELECT id FROM organization WHERE name = 'Test Organization'), (SELECT id FROM employee WHERE username = 'test_user'));
INSERT INTO tender (tender_id, tender_name, tender_description, service_type, status, organization_id, tender_version, created_at) VALUES
(uuid_generate_v4(), 'eu in', 'veniam adipisicing ea', 'Construction', 'Created', (SELECT id FROM organization WHERE name = 'Test Organization'), 1, CURRENT_TIMESTAMP);
