package com.github.halspals.smarthomeadapters.smarthomeadapters.model

import android.content.Context
import android.os.Handler
import android.util.Log
import android.view.MotionEvent
import android.widget.ImageView
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.EditRobotFragment
import com.github.halspals.smarthomeadapters.smarthomeadapters.MainActivity
import com.github.halspals.smarthomeadapters.smarthomeadapters.R
import com.github.halspals.smarthomeadapters.smarthomeadapters.RestApiService
import com.google.gson.annotations.SerializedName
import kotlinx.android.synthetic.main.activity_main.*
import net.openid.appauth.AuthorizationException
import okhttp3.ResponseBody
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response

data class RobotStatus(
        @SerializedName("value")
        var value: Boolean,
        @SerializedName("current")
        var current: Int,
        @SerializedName("min")
        val min: Int,
        @SerializedName("max")
        val max: Int
)

/**
 * Model for the body of a robot registration call
 *
 * @property nickname the nickname to set for the robot
 * @property robotType the type of functionality the robot provides
 */
data class RobotRegistrationBody(val nickname: String, val robotType: String)

/**
 * Model for a smart home adapter robot
 *
 * @property id the robot's unique ID
 * @property nickname name of the robot
 * @property robotType the type of functionality the robot provides (toggle, range, etc)
 * @property robotInterfaceType the type of interface/use case of the robot (switch, thermostat, etc)
 * @property robotStatus the robot's current state
 */
data class Robot(
    val id: String,
    val nickname: String,
    val robotType: String,
    @SerializedName("interfaceType") val robotInterfaceType: String,
    @SerializedName("status") val robotStatus: RobotStatus
) {
    companion object {
        const val INTERFACE_TYPE_TOGGLE = "toggle"
        const val INTERFACE_TYPE_RANGE = "range"
        const val ROBOT_TYPE_SWITCH = "switch"
        const val ROBOT_TYPE_THERMOSTAT = "thermostat"
        const val ROBOT_TYPE_BOLTLOCK = "boltlock"
        val ADD_ROBOT = Robot("","","","",
                RobotStatus(false,0,0,0))
    }

    private val tag = "RobotClass"

    internal fun updateViews(
            context: Context,
            robotCircle: ImageView? = null,
            robotIcon: ImageView? = null,
            robotRangeText: TextView? = null
    ) {
        when (robotInterfaceType) {
            Robot.INTERFACE_TYPE_TOGGLE -> {
                robotCircle?.setColorFilter(
                        if (robotStatus.value) {
                            context.getColor(R.color.colorToggleOn)
                        } else {
                            context.getColor(R.color.colorToggleOff)
                        }
                )
            }

            Robot.INTERFACE_TYPE_RANGE -> {
                robotRangeText?.text = robotStatus.current.toString()
            }

            else -> TODO("NO OTHER ROBOT INTERFACE TYPE EXPECTED")
        }

        when (robotType) {
            Robot.ROBOT_TYPE_SWITCH -> {
                robotIcon?.setImageResource(
                        if (robotStatus.value) {
                            R.drawable.ic_light_on
                        } else {
                            R.drawable.ic_light_off
                        }
                )
            }

            Robot.ROBOT_TYPE_THERMOSTAT -> {
                // Display no image
            }

            Robot.ROBOT_TYPE_BOLTLOCK -> {
                robotIcon?.setImageResource(
                        if (robotStatus.value) {
                            R.drawable.basic_lock
                        } else {
                            R.drawable.basic_lock_open
                        }
                )
            }

            else -> TODO("NO OTHER ROBOT TYPE EXPECTED")

        }
    }

    @SuppressWarnings("ClickableViewAccessibility")
    internal fun setViewEvents(
            parent: MainActivity,
            robotCircle: ImageView? = null,
            robotIcon: ImageView? = null,
            robotRangeText: TextView? = null
    ) {
        // configure interactions with the robot
        when (robotInterfaceType) {
            Robot.INTERFACE_TYPE_TOGGLE -> {
                // Set an onClick listener to handle boolean click events
                robotCircle?.setOnClickListener { _ ->
                    if (parent.isInEditMode) {
                        parent.robotToEdit = this
                        parent.startFragment(EditRobotFragment())
                    } else {
                        onToggle(parent, robotCircle, robotIcon)
                    }
                }
            }

            Robot.INTERFACE_TYPE_RANGE -> {
                // Set up the necessary objects to handle touch events
                val circleLocation = IntArray(2)
                val handler = Handler()
                val touchRunnableIncrease = object : Runnable {
                    override fun run() {
                        robotStatus.current++
                        if (robotStatus.current < robotStatus.max) {
                            robotStatus.current = robotStatus.max
                        }
                        Log.d(tag, "Increasing current value of $this")
                        updateViews(context = parent, robotRangeText = robotRangeText)
                        handler.postDelayed(this, 500)
                    }
                }
                val touchRunnableDecrease = object : Runnable {
                    override fun run() {
                        robotStatus.current--
                        if (robotStatus.current < robotStatus.min) {
                            robotStatus.current = robotStatus.min
                        }
                        Log.d(tag, "Increasing current value of $this")
                        updateViews(context = parent, robotRangeText = robotRangeText)
                        handler.postDelayed(this, 500)
                    }
                }

                // Set up touch events using these objects
                robotCircle?.setOnTouchListener { _, motionEvent ->
                    if (parent.isInEditMode) {
                        parent.robotToEdit = this
                        parent.startFragment(EditRobotFragment())
                    } else {
                        when (motionEvent.action) {
                            MotionEvent.ACTION_DOWN -> {
                                robotCircle.getLocationOnScreen(circleLocation)
                                if (motionEvent.rawY < circleLocation[1] + (robotCircle.height / 2)) {
                                    handler.post(touchRunnableIncrease)
                                } else {
                                    handler.post(touchRunnableDecrease)
                                }
                            }

                            MotionEvent.ACTION_UP -> {
                                handler.removeCallbacks(touchRunnableIncrease)
                                handler.removeCallbacks(touchRunnableDecrease)
                                onSeek(parent)
                            }

                            else -> {}
                        }
                    }
                    true
                }
            }
        }
    }

    /**
     * Handles a toggle-type event for a robot, updating its value and sending it to the server.
     */
    private fun onToggle(parent: MainActivity, robotCircle: ImageView?, robotIcon: ImageView?) {
        Log.wtf(tag, "[onToggle]: robot is $this)")

        // Send the update to the server
        parent.authState.performActionWithFreshTokens(parent.authService)
        { accessToken: String?, _: String?, ex: AuthorizationException? ->
            if (accessToken == null) {
                Log.e(tag, "[onSwitch] got null access token, exception: $ex")
            } else {
                parent.restApiService
                        .robotToggle(id, !robotStatus.value, accessToken, mapOf())
                        .enqueue(object: Callback<ResponseBody> {

                            override fun onResponse(
                                    call: Call<ResponseBody>,
                                    response: Response<ResponseBody>) {

                                if (response.isSuccessful) {
                                    parent.toast("Success")
                                    Log.d(tag, "[onToggle] Server accepted setting" +
                                            "toggle to ${!robotStatus.value}")
                                    robotStatus.value = !robotStatus.value
                                    updateViews(
                                            context = parent,
                                            robotCircle = robotCircle,
                                            robotIcon = robotIcon)

                                } else {
                                    val error = RestApiService.extractErrorFromResponse(response)
                                    Log.e(tag, "[onToggle] Unsuccessful, "
                                            + "error: $error")
                                    if (error != null) {
                                        parent.snackbar_layout.snackbar(error)
                                    }
                                }
                            }

                            override fun onFailure(call: Call<ResponseBody>, t: Throwable) {
                                val error = t.message
                                Log.e(tag, "[onToggle] FAILED, error: $error")
                                if (error != null) {
                                    parent.snackbar_layout.snackbar(error)
                                }
                            }
                        })
            }
        }
    }

    /**
     * onSeek is called whenever a range-type robot wants to make an API call to change state
     */
    private fun onSeek(parent: MainActivity) {
        Log.d(tag, "onSeek(${robotStatus.current})")

        parent.authState.performActionWithFreshTokens(parent.authService)
        { accessToken: String?, _: String?, ex: AuthorizationException? ->
            if (accessToken == null) {
                Log.e(tag, "[onSeek] got null access token, exception: $ex")
            } else {
                parent.restApiService
                        .robotRange(id, robotStatus.current, accessToken, mapOf())
                        .enqueue(object : Callback<ResponseBody> {

                            override fun onResponse(
                                    call: Call<ResponseBody>,
                                    response: Response<ResponseBody>) {

                                if (response.isSuccessful) {
                                    parent.toast("Success")
                                    Log.d(tag, "Server accepted setting range to ${robotStatus.current}")
                                } else {
                                    val error = RestApiService.extractErrorFromResponse(response)
                                    Log.e(tag, "Setting the range was unsuccessful, "
                                            + "error: $error")
                                    if (error != null) {
                                        parent.snackbar_layout.snackbar(error)
                                    }
                                }
                            }

                            override fun onFailure(call: Call<ResponseBody>, t: Throwable) {
                                val error = t.message
                                Log.e(tag, "onSeek(${robotStatus.current}) FAILED, error: $error")
                                if (error != null) {
                                    parent.snackbar_layout.snackbar(error)
                                }
                            }
                        })
            }
        }
    }
}
