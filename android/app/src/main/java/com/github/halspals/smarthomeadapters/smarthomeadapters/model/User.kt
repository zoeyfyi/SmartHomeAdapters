package com.github.halspals.smarthomeadapters.smarthomeadapters.model

data class User(val email: String, val password: String? = null) {
    override fun equals(other: Any?): Boolean {
        return if (other is User) {
            this.email == other.email
        } else {
            super.equals(other)
        }
    }

    override fun hashCode(): Int {
        return this.email.hashCode()
    }
}