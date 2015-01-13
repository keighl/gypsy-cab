r.dbDrop('gypsy_test');
r.dbCreate('gypsy_test');

r.db('gypsy_test').tableCreate('records');

r.db('gypsy_test').tableCreate('items');
r.db('gypsy_test').table('items').indexCreate('job_id');
r.db('gypsy_test').table('items').indexCreate('job_user_id');

r.db('gypsy_test').tableCreate('jobs');
r.db('gypsy_test').table('jobs').indexCreate('user_id');
r.db('gypsy_test').table('jobs').indexCreate('key');
r.db('gypsy_test').table('jobs').indexCreate('userAndKey', [r.row("user_id"), r.row("key")]);

r.db('gypsy_test').tableCreate('password_resets');
r.db('gypsy_test').table('password_resets').indexCreate('user_id');

r.db('gypsy_test').tableCreate('tokens');
r.db('gypsy_test').table('tokens').indexCreate('user_id');

r.db('gypsy_test').tableCreate('users');
r.db('gypsy_test').table('users').indexCreate('api_token');
r.db('gypsy_test').table('users').indexCreate('email');
