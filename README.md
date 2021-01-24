# udpecn
Showing how to send and receive explicit congestion notification using the UDP protocol.

This example does not do anything with the notification flags,
just shows how to set and receive them. You can run the example
on a single machine or on two machines. It uses ports 8001/8002
by default, sending an original message on 8001 and sending an
echo on 8002. Read the code or pass the -h flag to see other
options.
