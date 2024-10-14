/*
We should see how we handle the connection error problem...
1. If 1 source/replica couldn't connect, then the entire client will be declared as not workable
2. If any of source/replica doesn't work, then we will not declare as not workable until we have any left
   connections to connect to...
   We should also add here the reconnect thing...
   The reconnect should be explicitly used with care...
   Should we add here some kind of monitorer or retry process?!
   After we have connected, we use the PingFunction to check if it's alive or not...

I would like that the client will always check who's alive and who's not...
It will reconnect wherever possible... without disturbing anything...
Sometimes i don't care if it's connected or not... the resolver can handle that later...
We should also add to the connection pool somekind of Workable or not workable... meaning that is connected or not...
But all of them should be in the list from the beginning

i should make a lock on the resolver structure, so the manager of connections will add or notify that some connections
are not working....

everyone can have errors, but the policy should be simple and clear, without too much mechanisms which will handle them...
we should handle:
- on connect, if fail,  then retry a couple of times or infinitely
- on query -> if fail then go to next, one, if also failed circulate until is ok for a specific nr of times or seconds..
  if still failed then through to the user this error?!
  Every Request should have a response, we don't need the DB Client to handle all of the weight..., sometimes the user
  should do some work too.
- we should constantly check the dead connections and revive them!
- if replica failed, we should go to source
  if source failed, then retry...
- if everything failed, then the client totally failed... but, we should let the client live,
  we will still be accepting queries, but they should simply retry, and the monitorer should always check the other connections if they are up

Problems that we can encounter:
1. Initial Connection to the databases
2. Database failure...
3. Internet connection
4. Incorrect queries
5. invalid credentials
6. invalid driver...

*/











/*
When we have a transaction, we should memorize that has being started... committed or rollback
If we have some of these statements, then we should not switch connections, but we should check
the connections if they are working... meaning, before starting a transaction, we should check if the
connection is alive and is working... this means that it will handle the connections automatically...

usually when doing queries, we only know for registered callbacks, and we only know if they are inside of
a transaction... but we also need to know in what session the transaction has being started...

we have 2 statements:
PreparedStmtDB
and
PreparedStmtTX

so we should use somehow a specific source which will work well!
and we should also return error if something not works well

we should memorize that a transaction has being started? but there can be many...!!! and how we will
handle that? how we will know who called what?! do they have some mark that we can recognize them and attach
back to the same ConnPool?!
or we can simply check if it's still opened the transaction... if it's still opened, then we don't change
the connection pool, if it's not opened, then we choose one...!

Check the Session from GORM, that's where another session it's being created!
The session it's entirely new! it's a new *gorm.DB!!!

Everyone is calling Session to create a new instance... ->
methods like: callMethod


In conclusion

Each time when we do a request, we create a new gorm.DB Instance...
So even if we use Transaction, it will create a new one...
Now, the transaction is bonded to a single connection after it has started...
It's important for us to check the connection if it's ok before starting a transaction!

Use PingContext


So Transaction is called as a block... it has a fc -> function callback which receives the DB
so, we should check and know if it's creating new sessions or something like that...!


I should check if db.clone is higher than 0...
See Select function and getInstance which decides what to do...

because all statements are cloned between together... we can see if it's a transaction or not...


So we need to :
- know some kind id from the transaction
- know where it connects to...
- if the connection is dead... we should not switch the source connection, just restore it...
- when it finishes, it will not be a transaction anymore...
- we should make logs when testing...


- if the transaction has begun, then we should have the parameters or DSN of it... so that will be the destination,
  we also will need the connection id or session id! we can try using connection_id(), so it's unique based on DSN!
  ConnectionID() can be queried...


SO: MySQL does not support multiple sessions over a single connection, so 1 connection means 1 session!
Other drivers maybe support, but anyway, we should stick to 1 session 1 connection
Ok, but what if it's used by other goroutines at the same time, will they get into the transaction,
or another connection will be opened for them.

How can we read interface Properties... i think only through reflector or Methods that are declared there,
So we need to modify the Gorm library, or try using reflector!

ok... how can we test it?!...
1. Create a DB on local pc
2. create a client using real parameters....
3. Use the client from this lib/driver
   */
