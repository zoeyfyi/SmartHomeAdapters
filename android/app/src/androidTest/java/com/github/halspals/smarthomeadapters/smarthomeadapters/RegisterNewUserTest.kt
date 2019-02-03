package com.github.halspals.smarthomeadapters.smarthomeadapters


import android.support.test.espresso.Espresso.onView
import android.support.test.espresso.action.ViewActions.*
import android.support.test.espresso.intent.Intents
import android.support.test.espresso.intent.matcher.IntentMatchers
import android.support.test.espresso.intent.rule.IntentsTestRule
import android.support.test.espresso.matcher.ViewMatchers.*
import android.support.test.filters.LargeTest
import android.support.test.runner.AndroidJUnit4
import org.hamcrest.Matchers.allOf
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith

@LargeTest
@RunWith(AndroidJUnit4::class)
class RegisterNewUserTest {

    @Rule
    @JvmField
    var mActivityTestRule = IntentsTestRule(AuthenticationActivity::class.java)

    /**
     * Tests registering a new user in an "incremental" fashion, where the user
     * faces input validation errors at each step and corrects them only one at a time.
     */
    @Test
    fun registerNewUserTest() {
        val appCompatButton = onView(
                allOf(withId(R.id.register_button), withText("Register"),
                        isDisplayed()))
        appCompatButton.perform(click())


        val appCompatButton2 = onView(
                allOf(withId(R.id.register_button),
                        isDisplayed()))
        appCompatButton2.perform(click())

        val textInputEditText = onView(
                allOf(withId(R.id.email_input),
                        isDisplayed()))
        textInputEditText.perform(replaceText("test2@test.test"))

        appCompatButton2.perform(click())

        val textInputEditText2 = onView(
                allOf(withId(R.id.password_input),
                        isDisplayed()))
        textInputEditText2.perform(replaceText("testtest\$111"))

        val textInputEditText3 = onView(
                allOf(withId(R.id.confirm_password_input),
                        isDisplayed()))
        textInputEditText3.perform(replaceText("testtest\$11"))


        appCompatButton2.perform(click())

        textInputEditText3.perform(replaceText("testtest\$111"))


        appCompatButton2.perform(click())

        Intents.intended(IntentMatchers.hasComponent(MainActivity::class.java.name))
    }
}
