import 'package:flutter/material.dart';

class DashboardScreen extends StatelessWidget {
  const DashboardScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Attendance Dashboard'),
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
                      'Welcome to Attendance Management',
                      style: TextStyle(
                        fontSize: 20,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    SizedBox(height: 8),
                    Text(
                      'Manage your attendance, payroll and leave requests',
                      style: TextStyle(
                        color: Colors.grey,
                      ),
                    ),
                  ],
                ),
              ),
            ),

            // Feature Cards
            const Text(
              'Quick Access',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                marginBottom: 16,
              ),
            ),

            // Organization Management Card
            _buildFeatureCard(
              context,
              title: 'Organization',
              icon: Icons.people,
              color: Colors.blue,
              onTap: () {
                // Navigate to organization screen
              },
            ),

            const SizedBox(height: 16),

            // Attendance Tracking Card
            _buildFeatureCard(
              context,
              title: 'Attendance',
              icon: Icons.access_time,
              color: Colors.green,
              onTap: () {
                // Navigate to attendance screen
              },
            ),

            const SizedBox(height: 16),

            // Payroll Processing Card
            _buildFeatureCard(
              context,
              title: 'Payroll',
              icon: Icons.monetization_on,
              color: Colors.purple,
              onTap: () {
                // Navigate to payroll screen
              },
            ),

            const SizedBox(height: 16),

            // Leave Management Card
            _buildFeatureCard(
              context,
              title: 'Leave',
              icon: Icons.calendar_today,
              color: Colors.orange,
              onTap: () {
                // Navigate to leave screen
              },
            ),

            const SizedBox(height: 30),

            // System Overview
            const Text(
              'System Overview',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                marginBottom: 16,
              ),
            ),

            _buildStatsCard(context),
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

  Widget _buildStatsCard(BuildContext context) {
    return Card(
      elevation: 4,
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _buildStatItem('Employees', '1,248'),
                _buildStatItem('Attendance', '98.7%'),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _buildStatItem('Pending Leaves', '24'),
                _buildStatItem('Payroll', '100%'),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStatItem(String label, String value) {
    return Column(
      children: [
        Text(
          value,
          style: const TextStyle(
            fontSize: 24,
            fontWeight: FontWeight.bold,
            color: Colors.blue,
          ),
        ),
        Text(
          label,
          style: const TextStyle(
            color: Colors.grey,
          ),
        ),
      ],
    );
  }
}