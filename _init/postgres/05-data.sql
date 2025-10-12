\c demo_db;

-- Data initializations

-- Insert sample customers
INSERT INTO demo.customers (customer_number, full_name, id_number, phone_number, email, address, date_of_birth) VALUES
('CUST0000001', 'Ahmad Hidayat', '3201012345670001', '081234567801', 'ahmad.hidayat@email.com', 'Jl. Merdeka No. 123, Jakarta', '1985-03-15'),
('CUST0000002', 'Siti Nurhaliza', '3201012345670002', '081234567802', 'siti.nurhaliza@email.com', 'Jl. Sudirman No. 456, Jakarta', '1990-07-22'),
('CUST0000003', 'Muhammad Rizki', '3201012345670003', '081234567803', 'muhammad.rizki@email.com', 'Jl. Gatot Subroto No. 789, Jakarta', '1988-11-10'),
('CUST0000004', 'Fatimah Zahra', '3201012345670004', '081234567804', 'fatimah.zahra@email.com', 'Jl. Thamrin No. 321, Jakarta', '1992-05-08'),
('CUST0000005', 'Abdullah Rahman', '3201012345670005', '081234567805', 'abdullah.rahman@email.com', 'Jl. Asia Afrika No. 654, Jakarta', '1987-09-25'),
('CUST0000006', 'Khadijah Sari', '3201012345670006', '081234567806', 'khadijah.sari@email.com', 'Jl. Diponegoro No. 987, Bandung', '1991-12-14'),
('CUST0000007', 'Umar Hasan', '3201012345670007', '081234567807', 'umar.hasan@email.com', 'Jl. Ahmad Yani No. 147, Surabaya', '1989-04-30'),
('CUST0000008', 'Aisyah Putri', '3201012345670008', '081234567808', 'aisyah.putri@email.com', 'Jl. Pahlawan No. 258, Medan', '1993-08-17'),
('CUST0000009', 'Ibrahim Malik', '3201012345670009', '081234567809', 'ibrahim.malik@email.com', 'Jl. Gajah Mada No. 369, Semarang', '1986-01-23'),
('CUST0000010', 'Maryam Nur', '3201012345670010', '081234567810', 'maryam.nur@email.com', 'Jl. Hayam Wuruk No. 741, Yogyakarta', '1994-06-19'),
('CUST0000011', 'Yusuf Ali', '3201012345670011', '081234567811', 'yusuf.ali@email.com', 'Jl. Pemuda No. 852, Makassar', '1990-10-05'),
('CUST0000012', 'Zainab Hakim', '3201012345670012', '081234567812', 'zainab.hakim@email.com', 'Jl. Veteran No. 963, Palembang', '1988-02-28'),
('CUST0000013', 'Hamza Faruq', '3201012345670013', '081234567813', 'hamza.faruq@email.com', 'Jl. Imam Bonjol No. 159, Malang', '1991-07-12'),
('CUST0000014', 'Hafsa Amina', '3201012345670014', '081234567814', 'hafsa.amina@email.com', 'Jl. Cendana No. 357, Denpasar', '1989-11-26'),
('CUST0000015', 'Zakariya Amin', '3201012345670015', '081234567815', 'zakariya.amin@email.com', 'Jl. Anggrek No. 486, Balikpapan', '1992-03-09'),
('CUST0000016', 'Ruqayyah Salma', '3201012345670016', '081234567816', 'ruqayyah.salma@email.com', 'Jl. Melati No. 597, Manado', '1987-12-31'),
('CUST0000017', 'Ismail Hadi', '3201012345670017', '081234567817', 'ismail.hadi@email.com', 'Jl. Mawar No. 624, Pontianak', '1993-05-18'),
('CUST0000018', 'Safiya Laila', '3201012345670018', '081234567818', 'safiya.laila@email.com', 'Jl. Kenanga No. 735, Pekanbaru', '1990-09-07'),
('CUST0000019', 'Bilal Rashid', '3201012345670019', '081234567819', 'bilal.rashid@email.com', 'Jl. Dahlia No. 846, Jambi', '1986-04-21'),
('CUST0000020', 'Sumaya Nadia', '3201012345670020', '081234567820', 'sumaya.nadia@email.com', 'Jl. Tulip No. 957, Banjarmasin', '1994-08-14');

-- Insert accounts with various Islamic banking account types
INSERT INTO demo.accounts (account_number, customer_id, account_type, account_status, balance, currency) VALUES
-- Wadiah accounts (safekeeping/savings)
('1001000000001', 1, 'WADIAH', 'ACTIVE', 15750000.00, 'IDR'),
('1001000000002', 2, 'WADIAH', 'ACTIVE', 8250000.00, 'IDR'),
('1001000000003', 3, 'WADIAH', 'ACTIVE', 22500000.00, 'IDR'),
('1001000000004', 4, 'WADIAH', 'ACTIVE', 5000000.00, 'IDR'),
('1001000000005', 5, 'WADIAH', 'ACTIVE', 31200000.00, 'IDR'),

-- Mudharabah accounts (profit-sharing investment)
('2001000000001', 6, 'MUDHARABAH', 'ACTIVE', 125000000.00, 'IDR'),
('2001000000002', 7, 'MUDHARABAH', 'ACTIVE', 75000000.00, 'IDR'),
('2001000000003', 8, 'MUDHARABAH', 'ACTIVE', 200000000.00, 'IDR'),
('2001000000004', 9, 'MUDHARABAH', 'ACTIVE', 50000000.00, 'IDR'),
('2001000000005', 10, 'MUDHARABAH', 'ACTIVE', 180000000.00, 'IDR'),

-- Qard accounts (benevolent loan)
('3001000000001', 11, 'QARD', 'ACTIVE', 2500000.00, 'IDR'),
('3001000000002', 12, 'QARD', 'ACTIVE', 1750000.00, 'IDR'),
('3001000000003', 13, 'QARD', 'ACTIVE', 3000000.00, 'IDR'),
('3001000000004', 14, 'QARD', 'ACTIVE', 1250000.00, 'IDR'),
('3001000000005', 15, 'QARD', 'ACTIVE', 4500000.00, 'IDR'),

-- Additional mixed accounts for more test data
('1001000000006', 16, 'WADIAH', 'ACTIVE', 12000000.00, 'IDR'),
('2001000000006', 17, 'MUDHARABAH', 'ACTIVE', 95000000.00, 'IDR'),
('3001000000006', 18, 'QARD', 'ACTIVE', 2000000.00, 'IDR'),
('1001000000007', 19, 'WADIAH', 'ACTIVE', 18500000.00, 'IDR'),
('2001000000007', 20, 'MUDHARABAH', 'ACTIVE', 165000000.00, 'IDR');

-- Insert some initial transaction logs for testing
INSERT INTO demo.transaction_log (account_id, transaction_type, amount, balance_before, balance_after, reference_number, description) VALUES
(1, 'BALANCE_INQUIRY', NULL, 15750000.00, 15750000.00, 'INQ2025001001', 'Balance inquiry via API'),
(2, 'BALANCE_INQUIRY', NULL, 8250000.00, 8250000.00, 'INQ2025001002', 'Balance inquiry via API'),
(3, 'BALANCE_INQUIRY', NULL, 22500000.00, 22500000.00, 'INQ2025001003', 'Balance inquiry via API'),
(6, 'BALANCE_INQUIRY', NULL, 125000000.00, 125000000.00, 'INQ2025001004', 'Balance inquiry via API'),
(11, 'BALANCE_INQUIRY', NULL, 2500000.00, 2500000.00, 'INQ2025001005', 'Balance inquiry via API');