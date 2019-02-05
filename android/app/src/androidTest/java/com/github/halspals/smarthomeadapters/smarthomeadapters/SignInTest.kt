package com.github.halspals.smarthomeadapters.smarthomeadapters


import android.support.test.espresso.Espresso.onView
import android.support.test.espresso.action.ViewActions.*
import android.support.test.espresso.intent.rule.IntentsTestRule
import android.support.test.espresso.matcher.ViewMatchers.*
import android.support.test.filters.LargeTest
import android.support.test.espresso.intent.Intents.intended
import android.support.test.espresso.intent.matcher.IntentMatchers.hasComponent
import android.support.test.runner.AndroidJUnit4
import org.hamcrest.Matchers.allOf
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith

@LargeTest
@RunWith(AndroidJUnit4::class)
class SignInTest {

    @Rule
    @JvmField
    var mActivityTestRule = IntentsTestRule(AuthenticationActivity::class.java)

    /**
     * Tests signing in from the [AuthenticationActivity].
     */
    @Test
    fun signInTest() {
        val textInputEditText = onView(
                allOf(withId(R.id.email_input),
                        isDisplayed()))
        textInputEditText.perform(replaceText("test@test.test"), closeSoftKeyboard())

        val textInputEditText2 = onView(
                allOf(withId(R.id.password_input),
                        isDisplayed()))
        textInputEditText2.perform(replaceText("Testtest$111"), closeSoftKeyboard())

        val appCompatButton = onView(
                allOf(withId(R.id.sign_in_button), withText("Sign in"),
                        isDisplayed()))
        appCompatButton.perform(click())

        intended(hasComponent(MainActivity::class.java.name))
    }

}
