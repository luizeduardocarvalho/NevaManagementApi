# 📧 Clerk Webhook Examples

## 🔐 Webhook Headers

Clerk sends these headers with every webhook:

```
Content-Type: application/json
svix-id: msg_2Lh9nLJhxg84p5PrjG97mGsaPQN
svix-timestamp: 1673554664
svix-signature: v1,g0hM9SsE+OTPJTGt/ymKZ12+PkrBjjPZYm6qYSn2bOA=
```

## 📝 Webhook Event Types

Our system handles these Clerk events:
- `user.created` - New user registration
- `user.updated` - User profile changes  
- `user.deleted` - User account deletion

---

## 1️⃣ User Created Event

**Event Type**: `user.created`

**Use Case**: When someone signs up for LabFlux

### Full Payload Example:
```json
{
  "type": "user.created",
  "object": "event",
  "data": {
    "id": "user_2Lh9nLJhxg84p5PrjG97mGsaPQN",
    "object": "user",
    "username": null,
    "first_name": "John",
    "last_name": "Researcher",
    "image_url": "https://img.clerk.com/eyJ0eXBlIjoicHJveHkiLCJzcmMiOiJodHRwczovL2ltYWdlcy5jbGVyay5kZXYvb2F1dGhfZ29vZ2xlL2ltZ18yTGg5bkxKaHhnODRwNVByakc5N21Hc2FQUE4uanBlZyIsInMiOiJbeyJvcCI6InJlc2l6ZSIsInciOjE2MH1dIiwiaWF0IjoxNjczNTU0NjY0fQ.abc123",
    "has_image": true,
    "primary_email_address_id": "idn_2Lh9nLJhxg84p5PrjG97mGsaPQN",
    "primary_phone_number_id": null,
    "primary_web3_wallet_id": null,
    "password_enabled": true,
    "two_factor_enabled": false,
    "totp_enabled": false,
    "backup_code_enabled": false,
    "email_addresses": [
      {
        "id": "idn_2Lh9nLJhxg84p5PrjG97mGsaPQN",
        "object": "email_address",
        "email_address": "john.researcher@university.edu",
        "verification": {
          "status": "verified",
          "strategy": "email_code",
          "attempts": null,
          "expire_at": null
        },
        "linked_to": []
      }
    ],
    "phone_numbers": [],
    "web3_wallets": [],
    "external_accounts": [
      {
        "id": "idn_2Lh9nLJhxg84p5PrjG97mGsaPQN",
        "object": "external_account",
        "provider": "google",
        "identification_id": "idn_2Lh9nLJhxg84p5PrjG97mGsaPQN",
        "provider_user_id": "1234567890",
        "approved_scopes": "openid email profile",
        "email_address": "john.researcher@university.edu",
        "first_name": "John",
        "last_name": "Researcher",
        "image_url": "https://lh3.googleusercontent.com/a/abc123",
        "username": null,
        "public_metadata": {},
        "label": null,
        "verification": {
          "status": "verified",
          "strategy": "oauth_google",
          "attempts": null,
          "expire_at": 1673555564643
        }
      }
    ],
    "created_at": 1673554664643,
    "updated_at": 1673554664643,
    "last_sign_in_at": 1673554664643,
    "banned": false,
    "locked": false,
    "lockout_expires_in_seconds": null,
    "verification_attempts_remaining": 3,
    "public_metadata": {},
    "private_metadata": {},
    "unsafe_metadata": {}
  }
}
```

### Minimal Required Fields:
```json
{
  "type": "user.created",
  "object": "event", 
  "data": {
    "id": "user_2Lh9nLJhxg84p5PrjG97mGsaPQN",
    "first_name": "John",
    "last_name": "Researcher",
    "email_addresses": [
      {
        "id": "idn_2Lh9nLJhxg84p5PrjG97mGsaPQN",
        "email_address": "john.researcher@university.edu",
        "verification": {
          "status": "verified"
        }
      }
    ],
    "created_at": 1673554664643,
    "updated_at": 1673554664643
  }
}
```

---

## 2️⃣ User Updated Event

**Event Type**: `user.updated`

**Use Case**: When user changes their profile (name, email, etc.)

### Example: Name Change
```json
{
  "type": "user.updated",
  "object": "event",
  "data": {
    "id": "user_2Lh9nLJhxg84p5PrjG97mGsaPQN",
    "first_name": "Jonathan", 
    "last_name": "Researcher-Smith",
    "email_addresses": [
      {
        "id": "idn_2Lh9nLJhxg84p5PrjG97mGsaPQN",
        "email_address": "john.researcher@university.edu",
        "verification": {
          "status": "verified"
        }
      }
    ],
    "created_at": 1673554664643,
    "updated_at": 1673558264643
  }
}
```

### Example: Email Change
```json
{
  "type": "user.updated", 
  "object": "event",
  "data": {
    "id": "user_2Lh9nLJhxg84p5PrjG97mGsaPQN",
    "first_name": "John",
    "last_name": "Researcher",
    "email_addresses": [
      {
        "id": "idn_2Lh9nLJhxg84p5PrjG97mGsaPQN",
        "email_address": "j.researcher@newuniversity.edu",
        "verification": {
          "status": "verified"
        }
      },
      {
        "id": "idn_old_email_id",
        "email_address": "john.researcher@university.edu", 
        "verification": {
          "status": "unverified"
        }
      }
    ],
    "created_at": 1673554664643,
    "updated_at": 1673558264643
  }
}
```

---

## 3️⃣ User Deleted Event

**Event Type**: `user.deleted`

**Use Case**: When user deletes their account

### Example:
```json
{
  "type": "user.deleted",
  "object": "event", 
  "data": {
    "id": "user_2Lh9nLJhxg84p5PrjG97mGsaPQN",
    "object": "user",
    "deleted": true
  }
}
```

---

## 🧪 Testing with Postman

### 1. Basic User Created Test:
```json
{
  "type": "user.created",
  "object": "event",
  "data": {
    "id": "user_test_12345",
    "first_name": "Test",
    "last_name": "User",
    "email_addresses": [
      {
        "id": "email_test_12345", 
        "email_address": "test@labflux.com",
        "verification": {
          "status": "verified"
        }
      }
    ],
    "created_at": 1673554664643,
    "updated_at": 1673554664643
  }
}
```

### 2. User with Multiple Emails:
```json
{
  "type": "user.created",
  "object": "event",
  "data": {
    "id": "user_multi_email",
    "first_name": "Multi",
    "last_name": "Email",
    "email_addresses": [
      {
        "email_address": "unverified@example.com",
        "verification": { "status": "unverified" }
      },
      {
        "email_address": "verified@labflux.com", 
        "verification": { "status": "verified" }
      }
    ],
    "created_at": 1673554664643,
    "updated_at": 1673554664643
  }
}
```

### 3. User Update Test:
```json
{
  "type": "user.updated",
  "object": "event",
  "data": {
    "id": "user_test_12345",
    "first_name": "Updated",
    "last_name": "User", 
    "email_addresses": [
      {
        "email_address": "updated@labflux.com",
        "verification": { "status": "verified" }
      }
    ],
    "created_at": 1673554664643,
    "updated_at": 1673558264643
  }
}
```

### 4. User Deletion Test:
```json
{
  "type": "user.deleted",
  "object": "event",
  "data": {
    "id": "user_test_12345",
    "deleted": true
  }
}
```

---

## 🔧 Required Headers for Testing

When testing with Postman or curl, include these headers:

```
Content-Type: application/json
svix-id: msg_test_123
svix-timestamp: 1673554664
svix-signature: v1,test_signature_here
```

**Note**: Signature verification is currently disabled in development (TODO in code).

---

## 📊 Expected Responses

### ✅ Success (All Events):
```
HTTP 200 OK
```

### ❌ Errors:

**Invalid JSON:**
```json
{
  "error": "Invalid JSON payload"
}
```

**No Verified Email:**
```json
{
  "error": "No verified email found"
}
```

**User Not Found (for updates):**
```
HTTP 200 OK
(Event ignored, logged)
```

---

## 🔄 Integration Flow

1. **User signs up** → Clerk sends `user.created` → LabFlux creates user record
2. **User updates profile** → Clerk sends `user.updated` → LabFlux updates user
3. **User accepts invitation** → Uses `/api/auth/invite/{token}` → Links user to lab
4. **User deletes account** → Clerk sends `user.deleted` → LabFlux soft-deletes user

This webhook integration keeps LabFlux user data synchronized with Clerk authentication!