import 'package:flutter/material.dart';
import 'punch_screen.dart';

class EmployeePortalScreen extends StatelessWidget {
  const EmployeePortalScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Employee Portal'),
        backgroundColor: Colors.blue,
        foregroundColor: Colors.white,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Welcome Section
            const Card(
              elevation: 4,
              margin: EdgeInsets.only(bottom: 20),
              child: Padding(
                padding: EdgeInsets.all(16.0),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Welcome, John Smith',
                      style: TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    SizedBox(height: 8),
                    Text(
                      'Employee ID: EMP-12345',
                      style: TextStyle(
                        color: Colors.grey,
                      ),
                    ),
                  ],
                ),
              ),
            ),

            // Quick Actions
            const Text(
              'Quick Actions',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                marginBottom: 16,
              ),
            ),

            // Punch Attendance Card
            _buildFeatureCard(
              context,
              title: 'Punch Attendance',
              icon: Icons.access_time,
              color: Colors.blue,
              onTap: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => const PunchScreen(),
                  ),
                );
              },
            ),

            const SizedBox(height: 16),

            // View Leave Balance Card
            _buildFeatureCard(
              context,
              title: 'Leave Balance',
              icon: Icons.calendar_today,
              color: Colors.green,
              onTap: () {
                // Navigate to leave balance screen
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Leave balance view')),
                );
              },
            ),

            const SizedBox(height: 16),

            // View Salary Slip Card
            _buildFeatureCard(
              context,
              title: 'Salary Slip',
              icon: Icons.receipt,
              color: Colors.purple,
              onTap: () {
                // Navigate to salary slip screen
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('Salary slip view')),
                );
              },
            ),

            const SizedBox(height: 30),

            // Recent Activity
            const Text(
              'Recent Activity',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                marginBottom: 16,
              ),
            ),

            _buildActivityCard(context, 'Punched in', '9:05 AM today'),
            const SizedBox(height: 16),
            _buildActivityCard(context, 'Leave request submitted', '2 days ago'),
            const SizedBox(height: 16),
            _buildActivityCard(context, 'Salary slip generated', '1 week ago'),
          ],
        ),
      ),
    );
  }

  Widget _buildFeatureCard(BuildContext context, 
      {required String title, 
      required IconData icon, 
      required Color color,
      required VoidCallback onTap}) {
    return Card(
      elevation: 4,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Row(
            children: [
              Icon(icon, size: 40, color: color),
              const SizedBox(width: 16),
              Expanded(
                child: Text(
                  title,
                  style: const TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              const Icon(Icons.arrow_forward_ios),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildActivityCard(BuildContext context, String title, String time) {
    return Card(
      elevation: 2,
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Row(
          children: [
            Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: Colors.blue[100],
                borderRadius: BorderRadius.circular(20),
              ),
              child: const Icon(
                Icons.access_time,
                color: Colors.blue,
              ),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    title,
                    style: const TextStyle(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  Text(
                    time,
                    style: TextStyle(
                      color: Colors.grey[600],
                      fontSize: 14,
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}