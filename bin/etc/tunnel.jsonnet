local Options = {
  // Tunnel server URL
  URL: "https://127.0.0.1:9000/dev/api/v1/echo",
  // Do not verify tls certificate
  InsecureSkipVerify: true,
  // Use quic to connect to the server URL
  Quic: true,
};
local to="127.0.0.1";
[
  Options {
    // Listen on local address
    From: ':29444',
    // Forwarded server connection address
    To: 'tcp://'+to+':9444',
  },
  Options {
    From: ':20000',
    To: 'tcp://'+to+':10000',
  },
  Options {
    From: ':20001',
    To: 'tcp://'+to+':10001',
  },
  Options {
    From: ':20002',
    To: 'tcp://'+to+':10002',
  },
  Options {
    From: ':20050',
    To: 'tcp://'+to+':10050',
  },
]