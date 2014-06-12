require 'mina/git'

# Basic settings:
#   domain       - The hostname to SSH to.
#   deploy_to    - Path to deploy into.
#   repository   - Git repo to clone from. (needed by mina/git)
#   branch       - Branch name to deploy. (needed by mina/git)

set :domain, ENV["MINA_DOMAIN"]
set :deploy_to, '/var/www/upi'
set :repository, ENV["MINA_GIT"]
set :branch, 'master'

# They will be linked in the 'deploy:link_shared_paths' step.
set :shared_paths, ['.godeps', '.envrc', 'logs', 'uploads']

# Optional settings:
#   set :user, 'foobar'    # Username in the server to SSH to.
#   set :port, '30000'     # SSH port number.

# This task is the environment that is loaded for most commands, such as
# `mina deploy` or `mina rake`.
task :environment do
  queue! %[source .envrc]
end

desc "Install external dependencies"
task :'gpm:install' => :environment do
  queue! %[gpm install]
end

desc "Compile application"
task :'app:compile' => :environment do
  queue! %[go build -o upi]
end

# Put any custom mkdir's in here for when `mina setup` is ran.
# For Rails apps, we'll make some of the shared paths that are shared between
# all releases.
task :setup do
  queue! %[mkdir -p "#{deploy_to}/shared/.godeps"]
  queue! %[chmod g+rx,u+rwx "#{deploy_to}/shared/.godeps"]
  queue! %[mkdir -p "#{deploy_to}/shared/logs"]
  queue! %[chmod g+rx,u+rwx "#{deploy_to}/shared/logs"]
  queue! %[mkdir -p "#{deploy_to}/shared/uploads"]
  queue! %[chmod g+rx,u+rwx "#{deploy_to}/shared/uploads"]
  queue! %[touch #{deploy_to}/shared/.envrc]
end

desc "Deploys the current version to the server."
task :deploy do
  deploy do
    # Put things that will set up an empty directory into a fully set-up
    # instance of your project.
    invoke :'git:clone'
    invoke :'deploy:link_shared_paths'
    invoke :'gpm:install'
    invoke :'app:compile'

    to :launch do
      queue! %[touch #{deploy_to}/tmp/restart.txt]
    end
  end
end
