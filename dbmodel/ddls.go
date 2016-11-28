package dbmodel

// this file is used by an "init" application to
// create DB and create tables
// ideally this must be run only once
// this is used only to facilitate automatic
// creation of database by a go application.

//TableCreateSQL - create statements
var TableCreateSQL = []string{
	`CREATE DATABASE IF NOT EXISTS relive;`,

	`USE relive;`,

	`CREATE TABLE IF NOT EXISTS Account (
		  ID int(11) NOT NULL AUTO_INCREMENT,
		  PID int(11) NOT NULL,
          UserName VARCHAR(100) NOT NULL,
		  FirstName varchar(100) NOT NULL,
		  LastName varchar(100) DEFAULT NULL,
		  EmailID varchar(100) NOT NULL,
		  PasswdDigest varchar(512) NOT NULL,
          Salt VARCHAR(128) NOT NULL,
		  Role tinyint NOT NULL,
		  PRIMARY KEY (ID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 ;`,

	`CREATE TABLE IF NOT EXISTS Product (
		  ProductID int(11) NOT NULL AUTO_INCREMENT,
		  ProductType varchar(100) NOT NULL,
		  StoreSize int(11) NOT NULL,
		  Duration int(11) NOT NULL,
		  Amount int(11) NOT NULL,
		  PRIMARY KEY (ProductID)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 ;`,

	`CREATE TABLE IF NOT EXISTS Payment (
		  ID int(11) NOT NULL,
		  CCNumber varchar(100) NOT NULL,
		  BillingAddress varchar(100) NOT NULL,
		  CCExpiry varchar(100) NOT NULL,
		  CVVCode int(11) NOT NULL,
		  CONSTRAINT Payment_ibfk_1 FOREIGN KEY (ID) REFERENCES Account (ID) ON DELETE CASCADE ON UPDATE CASCADE 
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 ;`,

	`CREATE TABLE IF NOT EXISTS PaymentHistory (
		  ID int(11) NOT NULL,
		  LastPaidState varchar(100) NOT NULL,
		  LastType varchar(100) NOT NULL,
		  CONSTRAINT PaymentHistory_ibfk_1 FOREIGN KEY (ID) REFERENCES Account (ID) ON DELETE CASCADE ON UPDATE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 ;`,

	`CREATE TABLE IF NOT EXISTS Subscription (
		  ID int(11) NOT NULL,
		  ProductID int(11) NOT NULL,
		  ProductType varchar(100) NOT NULL,
		  StoreLocation varchar(100) NOT NULL,
		  StartDate TIMESTAMP DEFAULT '1970-01-01 00:00:01',
		  EndDate TIMESTAMP DEFAULT '1970-01-01 00:00:01',
		  NumberOfAdmins int(11) NOT NULL,
		  CONSTRAINT Subscription_ibfk_1 FOREIGN KEY (ID) REFERENCES Account (ID) ON DELETE CASCADE ON UPDATE CASCADE,
		  CONSTRAINT Subscription_ibfk_2 FOREIGN KEY (ProductID) REFERENCES Product (ProductID) ON DELETE CASCADE ON UPDATE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 ;`,

	`CREATE TABLE IF NOT EXISTS SubscriptionAccount (
		  ID int(11) NOT NULL,
		  PID int(11) NOT NULL,
		  CONSTRAINT SubscriptionAccount_ibfk_1 FOREIGN KEY (ID) REFERENCES Account (ID) ON DELETE CASCADE ON UPDATE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 ;`,

	`CREATE TABLE IF NOT EXISTS MediaType (
		  ID int(11) NOT NULL,
		  Catalog varchar(100) NOT NULL,
		  FileName varchar(100) DEFAULT NULL,
		  Title varchar(100) NOT NULL,
		  Description varchar(32) NOT NULL,
		  URL varchar(100) NOT NULL,
		  Poster varchar(10) NOT NULL,
		  CONSTRAINT MediaType_ibfk_1 FOREIGN KEY (ID) REFERENCES Account (ID) ON DELETE CASCADE ON UPDATE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 ;`,
}

//TableDeleteSQL - delete/drop statements
var TableDeleteSQL = []string{}
