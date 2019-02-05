package com.github.halspals.smarthomeadapters.smarthomeadapters


import android.support.test.espresso.Espresso.onView
import android.support.test.espresso.action.ViewActions.click
import android.support.test.espresso.assertion.ViewAssertions.matches
import android.support.test.espresso.matcher.ViewMatchers.*
import android.support.test.filters.LargeTest
import android.support.test.rule.ActivityTestRule
import android.support.test.runner.AndroidJUnit4
import org.hamcrest.Matchers.allOf
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith

@LargeTest
@RunWith(AndroidJUnit4::class)
class CycleLoginRegisterLoginTest {

    @Rule
    @JvmField
    var mActivityTestRule = ActivityTestRule(AuthenticationActivity::class.java)

    /**
     * Tests cycling from [AuthenticationActivity] to [RegisterUserActivity] and back again.
     */
    @Test
    fun cycleLoginRegisterLoginTest() {
        val appCompatButton = onView(
                allOf(withId(R.id.register_button), withText("Register"),
                        isDisplayed()))
        appCompatButton.perform(click())

        val appCompatButton2 = onView(
                allOf(withId(R.id.login_button), withText("Sign in instead"),
                        isDisplayed()))
        appCompatButton2.perform(click())

        val button = onView(
                allOf(withId(R.id.sign_in_button),
                        isDisplayed()))
        button.check(matches(isDisplayed()))
    }

}
