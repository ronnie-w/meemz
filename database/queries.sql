CREATE TABLE users(
    id int not null PRIMARY key AUTO_INCREMENT,
    username blob not null,
    userId varchar(255) not null,
    email varchar(255) not null,
    password blob not null,
    serverCode varchar(255) not null,
    passCode varchar(255) not null DEFAULT "0",
    isVerified varchar(3) not null DEFAULT "No",
    profileImg varchar(255) not null DEFAULT "blank-profile-picture-973460_1280.png",
    bio blob not null DEFAULT "No bio found"
);

CREATE TABLE posts(
    id int not null PRIMARY key AUTO_INCREMENT,
    userId varchar(50) not null,
    imgName varchar(255) not null,
    labelOcr blob not null,
    logoOcr blob not null,
    faceOcr blob not null,
    landmarkOcr blob not null,
    textOcr blob not null,
    safeSearchOcr blob not null,
    possibleDuplicate varchar(3) not null DEFAULT "No",
    duplicateNum blob not null DEFAULT "0",
    tags blob not null DEFAULT "",
    pComment blob not null DEFAULT "",
    access blob DEFAULT "Public",
    uploadTime blob DEFAULT "N/A"
);

CREATE TABLE chatRooms(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    topic blob not null,
    title blob not null,
    topicProfile blob not null,
    authorizedNumber blob not null,
    chatRoomId blob not null,
    createdOn blob not null
);

CREATE TABLE chatRoomMembers(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    chatRoomId blob not null,
    memberStatus blob not null,
    joinedOn blob not null
);

CREATE TABLE activeChats(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    chatType blob not null,
    chat blob not null,
    chatDate blob not null,
    chatTime blob not null,
    chatRoomId blob not null
);

CREATE TABLE reactions(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    imgName blob not null,
    reactionType blob not null
);

CREATE TABLE reports(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    imgName blob not null,
    report blob not null
);

CREATE TABLE comments(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    imgName blob not null,
    comment blob not null
);

CREATE TABLE regommend(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    imgName blob not null,
    pc blob not null
);

CREATE TABLE subs(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    creatorId blob not null,
    subTime blob not null
);

CREATE TABLE notifications(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    senderId blob not null,
    notify blob not null,
    receiveTime blob not null
);