CREATE TABLE users(
    id int not null PRIMARY key AUTO_INCREMENT,
    username blob not null,
    userId varchar(255) not null,
    email blob not null,
    password blob not null,
    serverCode varchar(255) not null,
    joinDate varchar(255) not null,
    passCode varchar(255) not null DEFAULT "0",
    isVerified varchar(3) not null DEFAULT "No",
    profileImg varchar(255) not null DEFAULT "blank-profile-picture-973460_1280.png",
    bio blob not null DEFAULT "No bio found"
);

CREATE TABLE posts(
    id int not null PRIMARY key AUTO_INCREMENT,
    userId varchar(50) not null,
    fileName varchar(255) not null,
    labelOcr blob not null DEFAULT "N/A",
    logoOcr blob not null DEFAULT "N/A",
    faceOcr blob not null DEFAULT "N/A",
    landmarkOcr blob not null DEFAULT "N/A",
    textOcr blob not null DEFAULT "N/A",
    safeSearchOcr blob not null DEFAULT "N/A",
    possibleDuplicate varchar(3) not null DEFAULT "No",
    duplicateNum blob not null DEFAULT "0",
    tags blob not null DEFAULT "",
    pComment blob not null DEFAULT "",
    credits blob not null DEFAULT "",
    uploadTime blob not null DEFAULT "N/A",
    originalName blob not null DEFAULT "",
    fileId blob not null DEFAULT "",
    fileIndex int not null DEFAULT 0
);

CREATE TABLE reactions(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    fileId blob not null,
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
    fileId blob not null,
    commentId blob not null,
    commentTime blob not null,
    comment blob not null
);

CREATE TABLE commentReplyLikes(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    fileId blob not null,
    commentReplyId blob not null
);

CREATE TABLE replies(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    fileId blob not null,
    commentId blob not null,
    replyId blob not null,
    replyTime blob not null,
    reply blob not null
);

CREATE TABLE regommend(
    id int not null PRIMARY KEY AUTO_INCREMENT,
    userId blob not null,
    imgName blob not null,
    pc blob not null,
    content_type blob not null 
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
    receiverId blob not null,
    notify blob not null,
    receiveTime blob not null
);