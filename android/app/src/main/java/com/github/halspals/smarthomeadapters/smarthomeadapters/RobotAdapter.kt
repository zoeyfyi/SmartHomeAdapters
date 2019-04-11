package com.github.halspals.smarthomeadapters.smarthomeadapters

import android.content.Context
import android.util.Log
import android.view.LayoutInflater
import android.view.MotionEvent
import android.view.View
import android.view.ViewGroup
import android.widget.BaseAdapter
import android.widget.ImageView
import android.widget.TextView
import com.github.halspals.smarthomeadapters.smarthomeadapters.model.Robot
import kotlinx.android.synthetic.main.activity_main.*
import net.openid.appauth.AuthorizationException
import okhttp3.ResponseBody
import org.jetbrains.anko.design.snackbar
import org.jetbrains.anko.startActivity
import org.jetbrains.anko.toast
import retrofit2.Call
import retrofit2.Callback
import retrofit2.Response


/**
 * Provides means of listing a series of [Robot]s with various functionality.
 */
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

    /**
     * Inflates and sets up the view corresponding to a non-dummy [Robot] register to the user.
     */
    private fun getRealRobotView(position: Int, convertView: View?, viewGroup: ViewGroup?): View {

        val robot = robots[position]

        // inflate card view
        val inflater = parent.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater
        val view: View = convertView ?: inflater.inflate(R.layout.view_robot_card, viewGroup, false)

        // get internal views
        val robotNickname = view.findViewById<TextView>(R.id.robot_nickname_text_view)
        val robotCircle = view.findViewById<ImageView>(R.id.robot_circle_drawable)
        val robotIcon = view.findViewById<ImageView>(R.id.robot_image_view)
        val robotRangeText = view.findViewById<TextView>(R.id.robot_range_text_view)

        robot.updateViews(parent, robotCircle, robotIcon, robotRangeText)
        robot.setViewEvents(parent, robotCircle, robotIcon, robotRangeText)



        // Configure static nickname view
        robotNickname.text = robots[position].nickname

        return view
    }

    /**
     * Inflates and sets up the view for the [Robot.ADD_ROBOT] dummy robot.
     */
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



}