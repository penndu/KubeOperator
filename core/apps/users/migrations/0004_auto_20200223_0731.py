# Generated by Django 2.2.10 on 2020-02-23 07:31

from django.db import migrations


class Migration(migrations.Migration):

    def forwards_func(apps, schema_editor):
        Profile = apps.get_model("users", "Profile")
        User = apps.get_model("auth", "User")
        db_alias = schema_editor.connection.alias
        admin = User.objects.using(db_alias).get(username='admin')
        if not admin.profile:
            Profile.objects.using(db_alias).create(user=admin)

    dependencies = [
        ('users', '0003_auto_20200223_0616'),
    ]

    operations = [
        migrations.RunPython(forwards_func),
    ]
