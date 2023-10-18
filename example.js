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
      console.log('attenuated token: ', attenuated);

      // inspect the attenuated token
      console.log('inspect : ', biscuit.inspect(attenuated));
    }
  catch(err){
      console.log(err.message);
  }
}
