/*Generated by orm 2022-01-31 00:56:48, from sqlite3 to mssql*/

IF NOT EXISTS (SELECT [name] FROM sys.tables WHERE [name] = '[test_dump_struct]' ) CREATE TABLE [test_dump_struct] ([id] INTEGER PRIMARY KEY IDENTITY NOT NULL, [name] VARCHAR(MAX) NULL, [is_man] INTEGER NULL, [created] DATETIME NULL);
SET IDENTITY_INSERT [test_dump_struct] ON;
INSERT INTO [test_dump_struct] ([id], [name], [is_man], [created]) VALUES (1,N'1',1,N'2022-01-30T19:26:48Z');
INSERT INTO [test_dump_struct] ([id], [name], [is_man], [created]) VALUES (2,N'2
',0,N'2022-01-30T19:26:48Z');
INSERT INTO [test_dump_struct] ([id], [name], [is_man], [created]) VALUES (3,N'3;',0,N'2022-01-30T19:26:48Z');
INSERT INTO [test_dump_struct] ([id], [name], [is_man], [created]) VALUES (4,N'4
;
''''',0,N'2022-01-30T19:26:48Z');
INSERT INTO [test_dump_struct] ([id], [name], [is_man], [created]) VALUES (5,N'5''
',0,N'2022-01-30T19:26:48Z');
