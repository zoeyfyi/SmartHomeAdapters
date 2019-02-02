package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
import android.util.Base64
import android.util.Log
import java.security.KeyStore
import javax.crypto.Cipher
import javax.crypto.KeyGenerator
import javax.crypto.SecretKey
import javax.crypto.spec.GCMParameterSpec

private const val ANDROID_KEY_STORE = "AndroidKeyStore"
private const val TRANSFORMATION = "AES/GCM/NoPadding"
private const val PREFERENCES_FILE = "TokenPreferences"
private const val TOKEN_KEY = "token"
private const val IV_KEY = "iv"

const val DEFAULT_ALIAS = "DoesThisNeedToBeSecret:thinkingface:"

private const val TAG = "KeyStoreManager"

/**
 * Stores the given token in a preference file, encrypting it and saving the key securely in Android KeyStore.
 *
 * @param token the authorization token to store
 * @param alias the alias to use for the KeyStore entry
 * @param context the context in which to get the shared preferences file
 */
internal fun securelyStoreToken(token: String, alias: String, context: Context) {

    // Set up an Android KeyStore key generator and get the secret key
    val keyGen = KeyGenerator.getInstance(KeyProperties.KEY_ALGORITHM_AES, ANDROID_KEY_STORE)

    val keyGenParamSpec = KeyGenParameterSpec.Builder(alias,
            KeyProperties.PURPOSE_ENCRYPT or KeyProperties.PURPOSE_DECRYPT)
            .setBlockModes(KeyProperties.BLOCK_MODE_GCM)
            .setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE)
            .build()

    keyGen.init(keyGenParamSpec)
    val secretKey: SecretKey = keyGen.generateKey()

    // Set up an encryption cipher and encrypt the token
    val cipher = Cipher.getInstance(TRANSFORMATION)
    cipher.init(Cipher.ENCRYPT_MODE, secretKey)
    val encrypted = cipher.doFinal(token.toByteArray())

    // Store the encrypted token along with the IV in shared prefs
    val sharedPrefs = context.getSharedPreferences(PREFERENCES_FILE, Context.MODE_PRIVATE)
    val editor = sharedPrefs.edit()
    editor.putString("$TOKEN_KEY|$alias", Base64.encodeToString(encrypted, Base64.DEFAULT))
    editor.putString("$IV_KEY|$alias", Base64.encodeToString(cipher.iv, Base64.DEFAULT))
    editor.apply()
}

/**
 * Retrieves a previously stored matching the KeyStore alias given.
 *
 * @param alias the KeyStore alias under which the keys were stored
 * @param context the context in which the shared prefs are located
 * @return the decrypted token, or null if there was no match in the prefs file
 */
internal fun retrieveStoredToken(alias: String, context: Context): String?  {

    // Retrieve the IV along with the encrypted token
    val sharedPrefs = context.getSharedPreferences(PREFERENCES_FILE, Context.MODE_PRIVATE)
    val ivString = sharedPrefs.getString("$IV_KEY|$alias", null)
    val encryptedTokenString = sharedPrefs.getString("$TOKEN_KEY|$alias", null)

    if (ivString == null) {
        Log.e(TAG, "IV not found for alias $alias")
        return null
    }
    if (encryptedTokenString == null) {
        Log.e(TAG, "Encrypted token not found for alias $alias")
        return null
    }

    // Decode them back to byte arrays
    val encryptedToken = Base64.decode(encryptedTokenString, Base64.DEFAULT)
    val iv = Base64.decode(ivString, Base64.DEFAULT)

    // Initialize a KeyStore instance and get the secret key for the alias
    val keyStore = KeyStore.getInstance(ANDROID_KEY_STORE)
    keyStore.load(null)
    val secretKey = (keyStore.getEntry(alias, null) as KeyStore.SecretKeyEntry).secretKey

    // Set up a GCM cipher to decipher the encrypted token with
    val spec = GCMParameterSpec(128, iv)
    val decipher = Cipher.getInstance(TRANSFORMATION)
    decipher.init(Cipher.DECRYPT_MODE, secretKey, spec)

    // Decrypt and return the token
    val decryptedToken = decipher.doFinal(encryptedToken)

    return String(decryptedToken)
}