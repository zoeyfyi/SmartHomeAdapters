package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.media.Image
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.BaseAdapter
import android.widget.ImageView
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import kotlinx.android.synthetic.main.fragment_robots.*
import net.openid.appauth.AuthorizationException
import okhttp3.ResponseBody
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.startActivity
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response


class RobotAdapter (
        private val parent: MainActivity,
        private val robots: MutableList<Robot>
) :  BaseAdapter() {

    private val tag = "RobotAdapter"

    init {
        // Add the dummy robot for starting the [RegisterRobotActivity]
        robots.add(Robot.ADD_ROBOT)
    }

    override fun getView(position: Int, convertView: View?, parent: ViewGroup?): View =
        if (position == count - 1) {
            assert(robots[position] == Robot.ADD_ROBOT) {
                "Getting view of last robot but it is not the ADD_ROBOT dummy; robot is " +
                        "${robots[position]}"
            }

            getAddRobotDummyView(convertView, parent)
        } else {
            getRealRobotView(position, convertView, parent)
    }

    private fun getRealRobotView(position: Int, convertView: View?, viewGroup: ViewGroup?): View {

        val robot = robots[position]

        // inflate card view
        val inflater = parent.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val view: View = convertView ?: inflater.inflate(R.layout.view_robot_card, viewGroup, false)

        // get internal views
        val robotNickname = view.findViewById<TextView>(R.id.robot_nickname_text_view)
        val robotCircle = view.findViewById<ImageView>(R.id.robot_circle_drawable)
        val robotIcon = view.findViewById<ImageView>(R.id.robot_image_view)

        // configure interactions with the robot
        when (robot.robotInterfaceType) {
            Robot.INTERFACE_TYPE_TOGGLE -> {
                robotCircle.setOnClickListener { _ -> onToggle(robot, robotCircle, robotIcon) }
            }

            Robot.INTERFACE_TYPE_RANGE -> {
                robotCircle.setOnDragListener { _, dragEvent ->
                    TODO("NOT IMPLEMENTED")
                }
            }
        }

        // Configure static nickname view
        robotNickname.text = robots[position].nickname

        return view
    }

    private fun getAddRobotDummyView(convertView: View?, viewGroup: ViewGroup?): View {
        // inflate card view
        val inflater = parent.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val view: View = convertView ?: inflater.inflate(R.layout.view_add_robot_card, viewGroup, false)

        // get internal views
        val robotCircle = view.findViewById<ImageView>(R.id.robot_circle_drawable)

        // configure click event for the dummy robot
        robotCircle.setOnClickListener { parent.startActivity<RegisterRobotActivity>() }

        return view
    }

    override fun getItem(position: Int): Any = robots[position]

    override fun getItemId(position: Int): Long = 0L

    override fun getCount(): Int = robots.size

    private fun updateRobotDisplay(robot: Robot, robotCircle: ImageView, robotIcon: ImageView) {
        when (robot.robotInterfaceType) {
            Robot.INTERFACE_TYPE_TOGGLE -> {
                robotIcon.setImageResource(R.drawable.basic_lightbulb)
                robotCircle.setColorFilter(
                        if (robot.robotStatus.value) {
                            parent.getColor(R.color.colorToggleOn)
                        } else {
                            parent.getColor(R.color.colorToggleOff)
                        }
                )
            }

            else -> TODO("NOT IMPLEMENTED")
        }
    }

    /**
     * Handles a toggle-type event for a robot, updating its value and sending it to the server.
     *
     * @param robot the [Robot] which the event fired for.
     */
    private fun onToggle(robot: Robot, robotCircle: ImageView, robotIcon: ImageView) {
        Log.d(tag, "[onToggle]: robot is $robot)")

        // Send the update to the server
        parent.authState.performActionWithFreshTokens(parent.authService)
        { accessToken: String?, _: String?, ex: AuthorizationException? ->
            if (accessToken == null) {
                Log.e(tag, "[onSwitch] got null access token, exception: $ex")
            } else {
                parent.restApiService
                        .robotToggle(robot.id, !robot.robotStatus.value, accessToken, mapOf())
                        .enqueue(object: Callback<ResponseBody> {

                            override fun onResponse(
                                    call: Call<ResponseBody>,
                                    response: Response<ResponseBody>) {

                                if (response.isSuccessful) {
                                    parent.toast("Success")
                                    Log.d(tag, "[onToggle] Server accepted setting" +
                                            "toggle to ${!robot.robotStatus.value}")
                                    robot.robotStatus.value = !robot.robotStatus.value
                                    updateRobotDisplay(robot, robotCircle, robotIcon)

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

}