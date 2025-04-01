create TABLE snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at DATETIME NOT NULL
);

create INDEX idx_snippets_created ON snippets (created_at);

insert into snippets (title, content, created_at, expired_at) values (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
    utc_timestamp(),
    date_add(utc_timestamp(), interval 365 day)
);
insert into snippets (title, content, created_at, expired_at) values (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    utc_timestamp(),
    date_add(utc_timestamp(), interval 365 day)
);
insert into snippets (title, content, created_at, expired_at) values (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    utc_timestamp(),
    date_add(utc_timestamp(), interval 7 day)
);
