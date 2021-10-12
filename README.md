# smsbroadcast
work sms broadcast program broke, bashed something up in go while we investigate to continue operations as normal

we have a CRM released in 1974 and a plugin for it written in house to make reqs to an sms broadcast api,
the only thing is we dont have the source code for the plugin and its written in .net 3.5 haha, anyways
this gets the same job done (but faster!!!)


update on the lore the .net 3.5 plugin uses tls 1.1 and the site we want to request uses tls 1.3
which is probably why it isnt very happy

i wish i just put all my vars in a struct instead of having 100 function arguments but its too late now
