-- Insert mock data into location table
INSERT INTO location (name, address, postal_code, city, country) VALUES
('Headquarters', '123 Main St', '10001', 'New York', 'USA'),
('Branch Office', '456 Side St', '20002', 'Washington', 'USA'),
('Remote Office', '789 Broad St', '30003', 'San Francisco', 'USA'),
('European Office', '12 Avenue des Champs', '75008', 'Paris', 'France'),
('Asia-Pacific Office', '34 Orchard Rd', '238832', 'Singapore', 'Singapore');

-- Insert mock data into department table
INSERT INTO department (name, location_id, manager_id) VALUES
('HR', 1, 1),
('Engineering', 1, 2),
('Sales', 2, 3),
('Marketing', 3, 4),
('Finance', 4, 5),
('IT', 5, 6);

-- Insert mock data into employee table
INSERT INTO employee (first_name, last_name, email, phone, hire_date, department_id, can_login, password) VALUES
('John', 'Doe', 'john.doe@example.com', '555-1234', '2020-01-15', 1, TRUE, '141780dc12e8a07c36e2cd28c975455b09328e6c65782c4152a1e18a4d802c98'), -- password: 'Admin01!'
('Jane', 'Smith', 'jane.smith@example.com', '555-5678', '2019-03-22', 2, FALSE, ''),
('Alice', 'Johnson', 'alice.johnson@example.com', '555-8765', '2018-07-30', 3, FALSE, ''),
('Bob', 'Brown', 'bob.brown@example.com', '555-4321', '2021-02-10', 4, FALSE, ''),
('Charlie', 'Davis', 'charlie.davis@example.com', '555-9101', '2017-11-05', 5, FALSE, ''),
('Diana', 'Garcia', 'diana.garcia@example.com', '555-1122', '2021-12-01', 6, FALSE, ''),
('Edward', 'Martinez', 'edward.martinez@example.com', '555-3344', '2020-08-18', 1, FALSE, ''),
('Fiona', 'Wilson', 'fiona.wilson@example.com', '555-5566', '2019-09-25', 2, FALSE, ''),
('George', 'Moore', 'george.moore@example.com', '555-7788', '2018-10-15', 3, FALSE, ''),
('Hannah', 'Taylor', 'hannah.taylor@example.com', '555-9900', '2017-12-20', 4, FALSE, ''),
('Ian', 'Anderson', 'ian.anderson@example.com', '555-2233', '2021-01-10', 5, FALSE, ''),
('Julia', 'Thomas', 'julia.thomas@example.com', '555-4455', '2020-05-12', 6, FALSE, ''),
('Kevin', 'Jackson', 'kevin.jackson@example.com', '555-6677', '2019-07-18', 1, FALSE, ''),
('Laura', 'White', 'laura.white@example.com', '555-8899', '2018-03-11', 2, FALSE, ''),
('Michael', 'Harris', 'michael.harris@example.com', '555-1010', '2017-08-07', 3, FALSE, ''),
('Nancy', 'Martin', 'nancy.martin@example.com', '555-1212', '2021-04-05', 4, FALSE, ''),
('Oliver', 'Lee', 'oliver.lee@example.com', '555-1313', '2020-06-17', 5, FALSE, ''),
('Patricia', 'Perez', 'patricia.perez@example.com', '555-1414', '2019-02-25', 6, FALSE, ''),
('Quincy', 'Thompson', 'quincy.thompson@example.com', '555-1515', '2018-11-14', 1, FALSE, ''),
('Rachel', 'Martinez', 'rachel.martinez@example.com', '555-1616', '2017-10-02', 2, FALSE, '');

