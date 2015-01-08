r.dbDrop('gypsy');
r.dbCreate('gypsy');

r.db('gypsy').tableCreate('items');
r.db('gypsy').table('items').indexCreate('job_id');
r.db('gypsy').table('items').indexCreate('job_user_id');

r.db('gypsy').tableCreate('jobs');
r.db('gypsy').table('jobs').indexCreate('user_id');
r.db('gypsy').table('jobs').indexCreate('key');
r.db('gypsy').table('jobs').indexCreate('userAndKey', [r.row("user_id"), r.row("key")]);

r.db('gypsy').tableCreate('password_resets');
r.db('gypsy').table('password_resets').indexCreate('token');

r.db('gypsy').tableCreate('tokens');
r.db('gypsy').table('tokens').indexCreate('token');
r.db('gypsy').table('tokens').indexCreate('user_id');

r.db('gypsy').tableCreate('users');
r.db('gypsy').table('users').indexCreate('api_token');
r.db('gypsy').table('users').indexCreate('email');
