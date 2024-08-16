DROP TABLE IF EXISTS record;
CREATE TABLE record (
  recordline INT AUTO_INCREMENT NOT NULL,
  recordsign VARCHAR(65,535) NOT NULL,
  recordcode VARCHAR(65,535) NOT NULL,
  recordmark VARCHAR(65,535) NOT NULL,
  recordhead VARCHAR(65,535),
  recordbody VARCHAR(65,535),
  recordtail VARCHAR(65,535),
  PRIMARY KEY (`recordline`)
);
DROP TABLE IF EXISTS member;
CREATE TABLE member (
  memberline INT AUTO_INCREMENT NOT NULL,
  memberpeer INT NOT NULL,
  membername VARCHAR(65,535) NOT NULL,
  PRIMARY KEY (`memberline`)
);
DROP TABLE IF EXISTS access;
CREATE TABLE access (
  accessnumber INT AUTO_INCREMENT NOT NULL,
  accessmember INT NOT NULL,
  accessrecord INT NOT NULL,
  PRIMARY KEY (`accessnumber`)
);