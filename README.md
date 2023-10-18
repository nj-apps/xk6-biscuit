# xk6-biscuit

This is a [k6](https://go.k6.io/k6) extension using the [xk6](https://github.com/grafana/xk6) system.

| :exclamation: This is a proof of concept, isn't supported by the k6 team, and may break in the future. USE AT YOUR OWN RISK! |
|------|

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Install `xk6`:
  ```shell
  $ go install go.k6.io/xk6/cmd/xk6@latest
  ```

2. Build the binary:
  ```shell
  $ xk6 build --with github.com/njapps/xk6-biscuit@latest
  ```

## Example

```javascript
import biscuit from 'k6/x/biscuit';

export const options = {
  scenarios: {
    sc1: {
      exec: 'attenuate_token',
      executor: 'per-vu-iterations',
      vus: 1,
    }
  },
};


export function attenuate_token() {
  // base64 encoded token 
  let token='EtQBCmoKDC9hL2ZpbGUxLnR4dAoML2EvZmlsZTIudHh0CgwvYS9maWxlMy50eHQSABgDIg0KCwgEEgMYgAgSAhgAIg0KCwgEEgMYgAgSAhgBIg0KCwgEEgMYgQgSAhgAIg0KCwgEEgMYgggSAhgBEiQIABIgkd8aEoGUBSHQ2xYBr8jS379Qc8Whf6xLyX321-hIMEAaQDwZUbmPBNOd2xMhH9cRS7-PTU7fI1MbuOmt47V0FHym7DM7gRFZerkdtuMwcRtQMWYKNmxh9MwN6GzNGNTB5gQiIgogbqzdI7E4NJueOH-XE6juHAFx5BrWq-52L4H_V1t45dw=';
  
  // Compute expiration date = now + 1h
  let expiration_date=new Date(Date.now());
  expiration_date.setHours(expiration_date.getHours()+1);

  // Array of blocks to append
  let blocks=['check if time($now), $now < '+expiration_date.toISOString(),
    'check if txn::service("Canal")'];

  try{
      // Attenuate the token
      let attenuated=biscuit.attenuate(token, blocks)
      console.log(attenuated);

      // inspect the attenuated token
      console.log('attenuated : ', biscuit.inspect(attenuated));
    }
  catch(err){
      console.log(err);
  }
}
```

Result output:

```
$ ./k6 run example.js
          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: example.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * sc1: 1 iterations for each of 1 VUs (maxDuration: 10m0s, exec: attenuate_token, gracefulStop: 30s)

INFO[0000] "2023-10-18T13:19:17.996Z"                    source=console
INFO[0000] object                                        source=console
INFO[0000] EtQBCmoKDC9hL2ZpbGUxLnR4dAoML2EvZmlsZTIudHh0CgwvYS9maWxlMy50eHQSABgDIg0KCwgEEgMYgAgSAhgAIg0KCwgEEgMYgAgSAhgBIg0KCwgEEgMYgQgSAhgAIg0KCwgEEgMYgggSAhgBEiQIABIgkd8aEoGUBSHQ2xYBr8jS379Qc8Whf6xLyX321-hIMEAaQDwZUbmPBNOd2xMhH9cRS7-PTU7fI1MbuOmt47V0FHym7DM7gRFZerkdtuMwcRtQMWYKNmxh9MwN6GzNGNTB5gQaxAEKWgoDbm93CgVDYW5hbAoMdHhuOjpzZXJ2aWNlEgAYAzIoCiYKAggbEgcIBRIDCIMIGhcKBQoDCIMICggKBiDl0r-pBgoEGgIIADIQCg4KAggbEggIhQgSAxiECBIkCAASICubspS9GkPRjMkSe5RvOn4-D3mRmjdy53LJfWmEQaQdGkCbDk3p8P-7obMQXVzqxD2ntDfDjK3Z8eftATxNHttD_QDcYqociiFpi6NsnG2pZYF9pf0b5tlbrGhXaSVxHnIPIiIKINWfJfSUR7JYOIFuPoxxHvts8Fk_ME8s30V3rCQYAxhJ  source=console
INFO[0000] attenuated :  
Biscuit {
	symbols: ["/a/file1.txt" "/a/file2.txt" "/a/file3.txt" "now" "Canal" "txn::service"]
	authority: Block {
		symbols: ["/a/file1.txt" "/a/file2.txt" "/a/file3.txt"]
		context: ""
		facts: [right("/a/file1.txt", "read") right("/a/file1.txt", "write") right("/a/file2.txt", "read") right("/a/file3.txt", "write")]
		rules: []
		checks: []
		version: 3
	}
	blocks: [Block {
		symbols: ["now" "Canal" "txn::service"]
		context: ""
		facts: []
		rules: []
		checks: [check if time($now), $now < 2023-10-18T14:19:17Z, check if txn::service("Canal")]
		version: 3
	}]
}  source=console

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=5.08ms min=5.08ms med=5.08ms max=5.08ms p(90)=5.08ms p(95)=5.08ms
     iterations...........: 1   192.380057/s


running (00m00.0s), 0/1 VUs, 1 complete and 0 interrupted iterations
sc1  ✓ [======================================] 1 VUs  00m00.0s/10m0s  1/1 iters, 1 per VU

```

Available functions :

* `inspect(b64Token string)` : deserialize the base64 token and returns token details as a string.
* `seal(b64Token string)` : prevent a token from being attenuated further. 'seal' returns a base64 encoded token.
* `attenuate(b64Token string, blocks []string)` : create a new token with the provided block appended. Return a base64 encoded token.
