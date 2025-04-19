# /// script
# dependencies = [
#   "dnslib",
# ]
# ///

"""
Simple DNS server for local testing/development.

Example /etc/resolver/dev.local:

nameserver 127.0.0.1
port 53
search_order 1
timeout 5
"""

from dnslib import DNSRecord, QTYPE, RR, A, DNSHeader
import socketserver

class DNSHandler(socketserver.BaseRequestHandler):
    def handle(self):
        data, sock = self.request
        request = DNSRecord.parse(data)
        qname = str(request.q.qname)
        reply = DNSRecord(DNSHeader(id=request.header.id, qr=1, aa=1, ra=1), q=request.q)
        if qname.endswith('.dev.local.'):
            reply.add_answer(RR(qname, QTYPE.A, rdata=A('127.0.0.1'), ttl=60))
        sock.sendto(reply.pack(), self.client_address)

if __name__ == '__main__':
    with socketserver.UDPServer(('0.0.0.0', 53), DNSHandler) as server:
        server.serve_forever()
